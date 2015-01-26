package spotify

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

// PlaylistTracks contains details about the tracks in a playlist.
type PlaylistTracks struct {
	// A link to the Web API endpoint where full details of
	// the playlist's tracks can be retrieved.
	Endpoint string `json:"href"`
	// The total number of tracks in the playlist.
	Total uint `json:"total"`
}

// SimplePlaylist contains basic info about a Spotify playlist.
type SimplePlaylist struct {
	// Indicates whether the playlist owner allows others to modify the playlist.
	// Note: only non-collaborative playlists are currently returned by Spotify's Web API.
	Collaborative bool `json:"collaborative"`
	// Known external URLs for this playlist.
	ExternalURLs ExternalURL `json:"external_urls"`
	// A link to the Web API endpoint providing full details of the playlist.
	Endpoint string `json:"href"`
	// The Spotify ID for the playlist.
	ID ID `json:"id"`
	// The playlist image.  Note: this field is only  returned for modified,
	// verified playlists. Otherwise the slice is empty.  If returned, the source
	// URL for the image is temporary and will expire in less than a day.
	Images []Image `json:"images"`
	// The name of the playlist.
	Name string `json:"name"`
	// The user who owns the playlist.
	Owner User `json:"owner"`
	// The playlist's public/private status
	IsPublic bool `json:"public"`
	// A collection to the Web API endpoint where full details of the playlist's
	// tracks can be retrieved, along with the total number of tracks in the playlist.
	Tracks PlaylistTracks `json:"tracks"`
	// The Spotify URI for the playlist.
	URI URI `json:"uri"`
}

// FullPlaylist provides extra playlist data in addition
// to the data provided by SimplePlaylist.
type FullPlaylist struct {
	SimplePlaylist
	// The playlist description.  Only returned for modified, verified playlists.
	Description string `json:"description"`
	// Information about the followers of this playlist.
	Followers Followers `json:"followers"`
	// The version identifier for the current playlist. Can be supplied in other
	// requests to target a specific playlist version.
	SnapshotID string `json:"snapshot_id"`
	// Information about the tracks of the playlist.
	// TODO: array of playlist track objects inside a
	// TODO: paging object.  is this the same as simple?
	Tracks string `json:"tracks"`
}

// PlaylistOptions contains optional parameters that can be used when querying
// for featured playlists.  Only the non-nil fields are used in the request.
type PlaylistOptions struct {
	Options
	// The desired language, consisting of a lowercase IO 639
	// language code and an uppercase ISO 3166-1 alpha-2
	// country code, joined by an underscore.  Provide this
	// parameter if you want the results returned in a particular
	// language.  If not specified, the result will be returned
	// in the Spotify default language (American English).
	Locale *string
	// A timestamp in ISO 8601 format (yyyy-MM-ddTHH:mm:ss).
	// use this paramter to specify th euser's local time to
	// get results tailored for that specific date and time
	// in the day.  If not provided, the response defaults to
	// the current UTC time.
	Timestamp *string
}

// FeaturedPlaylistsOpt gets a list of playlists featured by Spotify.
// It accepts a number of optional parameters via the opt argument.
func (c *Client) FeaturedPlaylistsOpt(opt *PlaylistOptions) (message string, playlists *PlaylistResult, e error) {
	if c.TokenType != BearerToken || c.AccessToken == "" {
		return "", nil, ErrAuthorizationRequired
	}
	uri := baseAddress + "browse/featured-playlists"
	if opt != nil {
		v := url.Values{}
		if opt.Locale != nil {
			v.Set("locale", *opt.Locale)
		}
		if opt.Country != nil {
			v.Set("country", *opt.Country)
		}
		if opt.Timestamp != nil {
			v.Set("timestamp", *opt.Timestamp)
		}
		if opt.Limit != nil {
			v.Set("limit", strconv.Itoa(*opt.Limit))
		}
		if opt.Offset != nil {
			v.Set("offset", strconv.Itoa(*opt.Offset))
		}
		if params := v.Encode(); params != "" {
			uri += "?" + params
		}
	}
	req, err := http.NewRequest("GET", uri, nil)
	if err != nil {
		return "", nil, errors.New("spotify: Couldn't create request")
	}
	req.Header.Set("Authorization", string(c.TokenType)+" "+c.AccessToken)
	resp, err := c.http.Do(req)
	if err != nil {
		return "", nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return "", nil, decodeError(resp.Body)
	}
	var result struct {
		Playlists *page  `json:"playlists"`
		Message   string `json:"message"`
	}
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return "", nil, err
	}
	return result.Message, toPlaylists(result.Playlists), nil
}

// FeaturedPlaylists gets a list of playlists featured by Spotify.
// It is equivalent to c.FeaturedPlaylistsOpt(nil).
func (c *Client) FeaturedPlaylists() (message string, playlists *PlaylistResult, e error) {
	return c.FeaturedPlaylistsOpt(nil)
}

// FollowPlaylist adds the current user as a follower of the specified
// playlist.  Any playlist can be followed, regardless of its private/public
// status, as long as you know the owner and playlist ID.
//
// If the public argument is true, then the playlist will be included in the
// user's public playlists.  To be able to follow playlists privately, the user
// must have granted the ScopePlaylistModifyPrivate scope.  The
// ScopePlaylistModifyPublic scope is required to follow playlists publicly.
func (c *Client) FollowPlaylist(owner ID, playlist ID, public bool) error {
	if c.TokenType != BearerToken || c.AccessToken == "" {
		return ErrAuthorizationRequired
	}
	uri := buildFollowURI(owner, playlist)
	body := strings.NewReader(strconv.FormatBool(public))
	req, err := http.NewRequest("PUT", uri, body)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", string(c.TokenType)+" "+c.AccessToken)
	// TODO: this is required - we should have a test to ensure it's in the header
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.http.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return decodeError(resp.Body)
	}
	return nil
}

// UnfollowPlaylist removes the current user as a follower of a playlist.
// This call requires authorization.  Unfollowing a publicly followed playlist
// requires the ScopePlaylistModifyPublic scope.  Unfolowing a privately followed,
// playlist requies the ScopePlaylistModifyPrivate scope.
func (c *Client) UnfollowPlaylist(owner, playlist ID) error {
	if c.TokenType != BearerToken || c.AccessToken == "" {
		return ErrAuthorizationRequired
	}
	uri := buildFollowURI(owner, playlist)
	req, err := http.NewRequest("DELETE", uri, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", string(c.TokenType)+" "+c.AccessToken)

	resp, err := c.http.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return decodeError(resp.Body)
	}
	return nil
}

func buildFollowURI(owner, playlist ID) string {
	buff := bytes.Buffer{}
	buff.WriteString(baseAddress)
	buff.WriteString("users/")
	buff.WriteString(string(owner))
	buff.WriteString("/playlists/")
	buff.WriteString(string(playlist))
	buff.WriteString("/followers")
	return string(buff.Bytes())
}
