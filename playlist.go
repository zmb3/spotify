// Copyright 2014, 2015 Zac Bergquist
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package spotify

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
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
	Tracks PlaylistTrackPage `json:"tracks"`
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
func (c *Client) FeaturedPlaylistsOpt(opt *PlaylistOptions) (message string, playlists *SimplePlaylistPage, e error) {
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
	req, err := c.newHTTPRequest("GET", uri, nil)
	if err != nil {
		return "", nil, errors.New("spotify: Couldn't create request")
	}
	resp, err := c.http.Do(req)
	if err != nil {
		return "", nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return "", nil, decodeError(resp.Body)
	}
	var result struct {
		Playlists SimplePlaylistPage `json:"playlists"`
		Message   string             `json:"message"`
	}
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return "", nil, err
	}
	return result.Message, &result.Playlists, nil
}

// FeaturedPlaylists gets a list of playlists featured by Spotify.
// It is equivalent to c.FeaturedPlaylistsOpt(nil).
func (c *Client) FeaturedPlaylists() (message string, playlists *SimplePlaylistPage, e error) {
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
	req, err := c.newHTTPRequest("PUT", uri, body)
	if err != nil {
		return err
	}
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
	req, err := c.newHTTPRequest("DELETE", uri, nil)
	if err != nil {
		return err
	}
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

// PlaylistsForUser gets a list of the playlists owned or followed by a particular
// Spotify user.  This call requires authorization.
//
// Private playlists are only retrievable for the current user, and require the
// ScopePlaylistReadPrivate scope.
//
// A user's collaborative playlists are not currently retrievable (this is a Web
// API limitation, not a limitation of package spotify).
func (c *Client) PlaylistsForUser(userID string) (*SimplePlaylistPage, error) {
	return c.PlaylistsForUserOpt(userID, nil)
}

// PlaylistsForUserOpt is like PlaylistsForUser, but it accepts optional paramters
// for filtering the results.
func (c *Client) PlaylistsForUserOpt(userID string, opt *Options) (*SimplePlaylistPage, error) {
	if c.TokenType != BearerToken || c.AccessToken == "" {
		return nil, ErrAuthorizationRequired
	}
	uri := baseAddress + "users/" + userID + "/playlists"
	if opt != nil {
		v := url.Values{}
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
	req, err := c.newHTTPRequest("GET", uri, nil)
	if err != nil {
		return nil, err
	}
	resp, err := c.http.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, decodeError(resp.Body)
	}
	var result SimplePlaylistPage
	err = json.NewDecoder(resp.Body).Decode(&result)
	return &result, err
}

// GetPlaylist gets a playlist owned by a Spotify user.
// This call requires authorization.  Both public and private
// playlists belonging to any user are retrievable with a valid
// access token.
func (c *Client) GetPlaylist(userID string, playlistID ID) (*FullPlaylist, error) {
	return c.GetPlaylistOpt(userID, playlistID, "")
}

// GetPlaylistOpt is like GetPlaylist, but it accepts an optional fields parameter
// that can be used to filter the query.
//
// fields is a comma-separated list of the fields to return.
// See the JSON tags on the FullPlaylist struct for valid field options.
// For example, to get just the playlist's description and URI:
//    fields = "description,uri"
//
// A dot separator can be used to specify non-reoccurring fields, while
// parentheses can be used to specify reoccurring fields within objects.
// For example, to get just the added date and the user ID of the adder:
//    fields = "tracks.items(added_at,added_by.id)"
//
// Use multiple parentheses to drill down into nested objects, for example:
//    fields = "tracks.items(track(name,href,album(name,href)))"
//
// Fields can be excluded by prefixing them with an exclamation mark, for example;
//    fields = "tracks.items(track(name,href,album(!name,href)))"
func (c *Client) GetPlaylistOpt(userID string, playlistID ID, fields string) (*FullPlaylist, error) {
	if c.TokenType != BearerToken || c.AccessToken == "" {
		return nil, ErrAuthorizationRequired
	}
	uri := baseAddress + "users/" + userID + "/playlists/" + string(playlistID)
	if fields != "" {
		uri += "?fields=" + url.QueryEscape(fields)
	}
	req, err := c.newHTTPRequest("GET", uri, nil)
	if err != nil {
		return nil, err
	}
	resp, err := c.http.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, decodeError(resp.Body)
	}
	var playlist FullPlaylist
	err = json.NewDecoder(resp.Body).Decode(&playlist)
	return &playlist, err
}

// GetPlaylistTracks gets full details of the tracks in a playlist, given the
// owner of the playlist and the playlist's Spotify ID.
// This call requires authorization.
func (c *Client) GetPlaylistTracks(userID string, playlistID ID) (*PlaylistTrackPage, error) {
	return c.GetPlaylistTracksOpt(userID, playlistID, nil, "")
}

// GetPlaylistTracksOpt is like GetPlaylistTracks, but it accepts optional parameters
// for sorting and filtering the results.  This call requries authorization.
//
// The field parameter is a comma-separated list of the fields to return.  See the
// JSON struct tags for the PlaylistTrackPage type for valid field names.
// For example, to get just the total number of tracks and the request limit:
//     fields = "total,limit"
//
// A dot separator can be used to specify non-reoccurring fields, while parentheses
// can be used to specify reoccurring fields within objects.  For example, to get
// just the added date and user ID of the adder:
//     fields = "items(added_at,added_by.id
//
// Use multiple parentheses to drill down into nested objects.  For example:
//     fields = "items(track(name,href,album(name,href)))"
//
// Fields can be excluded by prefixing them with an exclamation mark.  For example:
//     fields = "items.track.album(!external_urls,images)"
func (c *Client) GetPlaylistTracksOpt(userID string, playlistID ID, opt *Options, fields string) (*PlaylistTrackPage, error) {
	if c.TokenType != BearerToken || c.AccessToken == "" {
		return nil, ErrAuthorizationRequired
	}
	uri := fmt.Sprintf("%susers/%s/playlists/%s/tracks", baseAddress, userID, playlistID)
	v := url.Values{}
	if fields != "" {
		v.Set("fields", fields)
	}
	if opt != nil {
		if opt.Limit != nil {
			v.Set("limit", strconv.Itoa(*opt.Limit))
		}
		if opt.Offset != nil {
			v.Set("offset", strconv.Itoa(*opt.Offset))
		}
	}
	if params := v.Encode(); params != "" {
		uri += "?" + params
	}

	req, err := c.newHTTPRequest("GET", uri, nil)
	if err != nil {
		return nil, err
	}
	resp, err := c.http.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, decodeError(resp.Body)
	}
	var result PlaylistTrackPage
	err = json.NewDecoder(resp.Body).Decode(&result)
	return &result, err
}

// CreatePlaylistForUser creates a playlist for a Spotify user.
// The playlist will be empty until you add tracks to it.
// The playlistName does not need to be unique - a user can have
// several playlists with the same name.
//
// This call requires authorization.  Creating a public playlist
// for a user requires the ScopePlaylistModifyPublic scope;
// creating a private playlist requires the ScopePlaylistModifyPrivate
// scope.
//
// On success, the newly created playlist is returned.
func (c *Client) CreatePlaylistForUser(userID, playlistName string, public bool) (*FullPlaylist, error) {
	if c.TokenType != BearerToken || c.AccessToken == "" {
		return nil, ErrAuthorizationRequired
	}
	uri := fmt.Sprintf("%susers/%s/playlists", baseAddress, userID)
	body := struct {
		Name   string `json:"name"`
		Public bool   `json:"public"`
	}{
		playlistName,
		public,
	}
	bodyJSON, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	req, err := c.newHTTPRequest("POST", uri, bytes.NewReader(bodyJSON))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := c.http.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return nil, decodeError(resp.Body)
	}
	var p FullPlaylist
	err = json.NewDecoder(resp.Body).Decode(&p)
	return &p, err
}

// ChangePlaylistName changes the name of a playlist.  This call requires that the
// user has authorized the ScopePlaylistModifyPublic or ScopePlaylistModifyPrivate
// scopes (depending on whether the playlist is public or private).
// The current user must own the playlist in order to modify it.
func (c *Client) ChangePlaylistName(userID string, playlistID ID, newName string) error {
	return c.modifyPlaylist(userID, playlistID, newName, nil)
}

// ChangePlaylistAccess modifies the public/private status of a playlist.  This call
// requires that the user has authorized the ScopePlaylistModifyPublic or ScopePlaylistModifyPrivate
// scopes (depending on whether the playlist is currently public or private).
// The current user must own the playlist in order to modify it.
func (c *Client) ChangePlaylistAccess(userID string, playlistID ID, public bool) error {
	return c.modifyPlaylist(userID, playlistID, "", &public)
}

// ChangePlaylistNameAndAccess combines ChangePlaylistName and ChangePlaylistAccess into
// a single Web API call.  It requires that the user has authorized the ScopePlaylistModifyPublic
// or ScopePlaylistModifyPrivate scopes (depending on whether the playlist is currently
// public or private).  The current user must own the playlist in order to modify it.
func (c *Client) ChangePlaylistNameAndAccess(userID string, playlistID ID, newName string, public bool) error {
	return c.modifyPlaylist(userID, playlistID, newName, &public)
}

func (c *Client) modifyPlaylist(userID string, playlistID ID, newName string, public *bool) error {
	if c.TokenType != BearerToken || c.AccessToken == "" {
		return ErrAuthorizationRequired
	}
	body := struct {
		Name   string `json:"name,omitempty"`
		Public *bool  `json:"public,omitempty"`
	}{
		newName,
		public,
	}
	bodyJSON, err := json.Marshal(body)
	if err != nil {
		return err
	}
	uri := fmt.Sprintf("%susers/%s/playlists/%s", baseAddress, userID, string(playlistID))
	req, err := c.newHTTPRequest("PUT", uri, bytes.NewReader(bodyJSON))
	if err != nil {
		return err
	}
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

// AddTracksToPlaylist adds one or more tracks to a user's playlist.  This call requires
// authorization (ScopePlaylistModifyPublic or ScopePlaylistModifyPrivate).  A maximum of
// 100 tracks can be added per call.  It returns a snapshot ID that can be used to
// identify this version (the new version) of the playlist in future requests.
func (c *Client) AddTracksToPlaylist(userID string, playlistID ID, trackIDs ...ID) (snapshotID string, err error) {
	if c.TokenType != BearerToken || c.AccessToken == "" {
		return "", ErrAuthorizationRequired
	}
	// convert track IDs to Spotify URIs (spotify:track:<ID>)
	uris := make([]string, len(trackIDs))
	for i, id := range trackIDs {
		uris[i] = fmt.Sprintf("spotify:track:%s", id)
	}
	uri := fmt.Sprintf("%susers/%s/playlists/%s/tracks?uris=%s",
		baseAddress, userID, string(playlistID), strings.Join(uris, ","))
	req, err := c.newHTTPRequest("POST", uri, nil)
	if err != nil {
		return "", err
	}
	resp, err := c.http.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusCreated {
		return "", decodeError(resp.Body)
	}
	body := struct {
		SnapshotID string `json:"snapshot_id"`
	}{}
	err = json.NewDecoder(resp.Body).Decode(&body)
	if err != nil {
		// the response code indicates success..
		return "", err
	}
	return body.SnapshotID, nil
}
