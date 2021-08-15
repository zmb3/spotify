package spotify

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
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
	Collaborative bool              `json:"collaborative"`
	ExternalURLs  map[string]string `json:"external_urls"`
	// A link to the Web API endpoint providing full details of the playlist.
	Endpoint string `json:"href"`
	ID       ID     `json:"id"`
	// The playlist image.  Note: this field is only  returned for modified,
	// verified playlists. Otherwise the slice is empty.  If returned, the source
	// URL for the image is temporary and will expire in less than a day.
	Images   []Image `json:"images"`
	Name     string  `json:"name"`
	Owner    User    `json:"owner"`
	IsPublic bool    `json:"public"`
	// The version identifier for the current playlist. Can be supplied in other
	// requests to target a specific playlist version.
	SnapshotID string `json:"snapshot_id"`
	// A collection to the Web API endpoint where full details of the playlist's
	// tracks can be retrieved, along with the total number of tracks in the playlist.
	Tracks PlaylistTracks `json:"tracks"`
	URI    URI            `json:"uri"`
}

// FullPlaylist provides extra playlist data in addition to the data provided by SimplePlaylist.
type FullPlaylist struct {
	SimplePlaylist
	// The playlist description.  Only returned for modified, verified playlists.
	Description string `json:"description"`
	// Information about the followers of this playlist.
	Followers Followers         `json:"followers"`
	Tracks    PlaylistTrackPage `json:"tracks"`
}

// FeaturedPlaylistsOpt gets a list of playlists featured by Spotify.
// Supported options: Locale, Country, Timestamp, Limit, Offset
func (c *Client) FeaturedPlaylists(ctx context.Context, opts ...RequestOption) (message string, playlists *SimplePlaylistPage, e error) {
	spotifyURL := c.baseURL + "browse/featured-playlists"
	if params := processOptions(opts...).urlParams.Encode(); params != "" {
		spotifyURL += "?" + params
	}

	var result struct {
		Playlists SimplePlaylistPage `json:"playlists"`
		Message   string             `json:"message"`
	}

	err := c.get(ctx, spotifyURL, &result)
	if err != nil {
		return "", nil, err
	}

	return result.Message, &result.Playlists, nil
}

// FollowPlaylist adds the current user as a follower of the specified
// playlist.  Any playlist can be followed, regardless of its private/public
// status, as long as you know the playlist ID.
//
// If the public argument is true, then the playlist will be included in the
// user's public playlists.  To be able to follow playlists privately, the user
// must have granted the ScopePlaylistModifyPrivate scope.  The
// ScopePlaylistModifyPublic scope is required to follow playlists publicly.
func (c *Client) FollowPlaylist(ctx context.Context, playlist ID, public bool) error {
	spotifyURL := buildFollowURI(c.baseURL, playlist)
	body := strings.NewReader(strconv.FormatBool(public))
	req, err := http.NewRequestWithContext(ctx, "PUT", spotifyURL, body)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	err = c.execute(req, nil)
	if err != nil {
		return err
	}
	return nil
}

// UnfollowPlaylist removes the current user as a follower of a playlist.
// Unfollowing a publicly followed playlist requires ScopePlaylistModifyPublic.
// Unfolowing a privately followed playlist requies ScopePlaylistModifyPrivate.
func (c *Client) UnfollowPlaylist(ctx context.Context, playlist ID) error {
	spotifyURL := buildFollowURI(c.baseURL, playlist)
	req, err := http.NewRequestWithContext(ctx, "DELETE", spotifyURL, nil)
	if err != nil {
		return err
	}
	err = c.execute(req, nil)
	if err != nil {
		return err
	}
	return nil
}

func buildFollowURI(url string, playlist ID) string {
	return fmt.Sprintf("%splaylists/%s/followers",
		url, string(playlist))
}

// GetPlaylistsForUser gets a list of the playlists owned or followed by a
// particular Spotify user.
//
// Private playlists and collaborative playlists are only retrievable for the
// current user.  In order to read private playlists, the user must have granted
// the ScopePlaylistReadPrivate scope.  Note that this scope alone will not
// return collaborative playlists, even though they are always private.  In
// order to read collaborative playlists, the user must have granted the
// ScopePlaylistReadCollaborative scope.
//
// Supported options: Limit, Offset
func (c *Client) GetPlaylistsForUser(ctx context.Context, userID string, opts ...RequestOption) (*SimplePlaylistPage, error) {
	spotifyURL := c.baseURL + "users/" + userID + "/playlists"
	if params := processOptions(opts...).urlParams.Encode(); params != "" {
		spotifyURL += "?" + params
	}

	var result SimplePlaylistPage

	err := c.get(ctx, spotifyURL, &result)
	if err != nil {
		return nil, err
	}

	return &result, err
}

// GetPlaylist fetches a playlist from spotify.
// Supported options: Fields
func (c *Client) GetPlaylist(ctx context.Context, playlistID ID, opts ...RequestOption) (*FullPlaylist, error) {
	spotifyURL := fmt.Sprintf("%splaylists/%s", c.baseURL, playlistID)
	if params := processOptions(opts...).urlParams.Encode(); params != "" {
		spotifyURL += "?" + params
	}

	var playlist FullPlaylist

	err := c.get(ctx, spotifyURL, &playlist)
	if err != nil {
		return nil, err
	}

	return &playlist, err
}

// GetPlaylistTracks gets full details of the tracks in a playlist, given the
// playlist's Spotify ID.
//
// Supported options: Limit, Offset, Market, Fields
func (c *Client) GetPlaylistTracks(
	ctx context.Context,
	playlistID ID,
	opts ...RequestOption,
) (*PlaylistTrackPage, error) {
	spotifyURL := fmt.Sprintf("%splaylists/%s/tracks", c.baseURL, playlistID)
	if params := processOptions(opts...).urlParams.Encode(); params != "" {
		spotifyURL += "?" + params
	}

	var result PlaylistTrackPage

	err := c.get(ctx, spotifyURL, &result)
	if err != nil {
		return nil, err
	}

	return &result, err
}

// CreatePlaylistForUser creates a playlist for a Spotify user.
// The playlist will be empty until you add tracks to it.
// The playlistName does not need to be unique - a user can have
// several playlists with the same name.
//
// Creating a public playlist for a user requires ScopePlaylistModifyPublic;
// creating a private playlist requires ScopePlaylistModifyPrivate.
//
// On success, the newly created playlist is returned.
func (c *Client) CreatePlaylistForUser(ctx context.Context, userID, playlistName, description string, public bool, collaborative bool) (*FullPlaylist, error) {
	spotifyURL := fmt.Sprintf("%susers/%s/playlists", c.baseURL, userID)
	body := struct {
		Name          string `json:"name"`
		Public        bool   `json:"public"`
		Description   string `json:"description"`
		Collaborative bool   `json:"collaborative"`
	}{
		playlistName,
		public,
		description,
		collaborative,
	}
	bodyJSON, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequestWithContext(ctx, "POST", spotifyURL, bytes.NewReader(bodyJSON))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	var p FullPlaylist
	err = c.execute(req, &p, http.StatusCreated)
	if err != nil {
		return nil, err
	}

	return &p, err
}

// ChangePlaylistName changes the name of a playlist.  This call requires that the
// user has authorized the ScopePlaylistModifyPublic or ScopePlaylistModifyPrivate
// scopes (depending on whether the playlist is public or private).
// The current user must own the playlist in order to modify it.
func (c *Client) ChangePlaylistName(ctx context.Context, playlistID ID, newName string) error {
	return c.modifyPlaylist(ctx, playlistID, newName, "", nil)
}

// ChangePlaylistAccess modifies the public/private status of a playlist.  This call
// requires that the user has authorized the ScopePlaylistModifyPublic or
// ScopePlaylistModifyPrivate scopes (depending on whether the playlist is
// currently public or private).  The current user must own the playlist in order to modify it.
func (c *Client) ChangePlaylistAccess(ctx context.Context, playlistID ID, public bool) error {
	return c.modifyPlaylist(ctx, playlistID, "", "", &public)
}

// ChangePlaylistDescription modifies the description of a playlist.  This call
// requires that the user has authorized the ScopePlaylistModifyPublic or
// ScopePlaylistModifyPrivate scopes (depending on whether the playlist is
// currently public or private).  The current user must own the playlist in order to modify it.
func (c *Client) ChangePlaylistDescription(ctx context.Context, playlistID ID, newDescription string) error {
	return c.modifyPlaylist(ctx, playlistID, "", newDescription, nil)
}

// ChangePlaylistNameAndAccess combines ChangePlaylistName and ChangePlaylistAccess into
// a single Web API call.  It requires that the user has authorized the ScopePlaylistModifyPublic
// or ScopePlaylistModifyPrivate scopes (depending on whether the playlist is currently
// public or private).  The current user must own the playlist in order to modify it.
func (c *Client) ChangePlaylistNameAndAccess(ctx context.Context, playlistID ID, newName string, public bool) error {
	return c.modifyPlaylist(ctx, playlistID, newName, "", &public)
}

// ChangePlaylistNameAccessAndDescription combines ChangePlaylistName, ChangePlaylistAccess, and
// ChangePlaylistDescription into a single Web API call.  It requires that the user has authorized
// the ScopePlaylistModifyPublic or ScopePlaylistModifyPrivate scopes (depending on whether the
// playlist is currently public or private).  The current user must own the playlist in order to modify it.
func (c *Client) ChangePlaylistNameAccessAndDescription(ctx context.Context, playlistID ID, newName, newDescription string, public bool) error {
	return c.modifyPlaylist(ctx, playlistID, newName, newDescription, &public)
}

func (c *Client) modifyPlaylist(ctx context.Context, playlistID ID, newName, newDescription string, public *bool) error {
	body := struct {
		Name        string `json:"name,omitempty"`
		Public      *bool  `json:"public,omitempty"`
		Description string `json:"description,omitempty"`
	}{
		newName,
		public,
		newDescription,
	}
	bodyJSON, err := json.Marshal(body)
	if err != nil {
		return err
	}
	spotifyURL := fmt.Sprintf("%splaylists/%s", c.baseURL, string(playlistID))
	req, err := http.NewRequestWithContext(ctx, "PUT", spotifyURL, bytes.NewReader(bodyJSON))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	err = c.execute(req, nil, http.StatusCreated)
	if err != nil {
		return err
	}
	return nil
}

// AddTracksToPlaylist adds one or more tracks to a user's playlist.
// This call requires ScopePlaylistModifyPublic or ScopePlaylistModifyPrivate.
// A maximum of 100 tracks can be added per call.  It returns a snapshot ID that
// can be used to identify this version (the new version) of the playlist in
// future requests.
func (c *Client) AddTracksToPlaylist(ctx context.Context, playlistID ID, trackIDs ...ID) (snapshotID string, err error) {
	uris := make([]string, len(trackIDs))
	for i, id := range trackIDs {
		uris[i] = fmt.Sprintf("spotify:track:%s", id)
	}
	m := make(map[string]interface{})
	m["uris"] = uris

	spotifyURL := fmt.Sprintf("%splaylists/%s/tracks",
		c.baseURL, string(playlistID))
	body, err := json.Marshal(m)
	if err != nil {
		return "", err
	}
	req, err := http.NewRequestWithContext(ctx, "POST", spotifyURL, bytes.NewReader(body))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")

	result := struct {
		SnapshotID string `json:"snapshot_id"`
	}{}

	err = c.execute(req, &result, http.StatusCreated)
	if err != nil {
		return "", err
	}

	return result.SnapshotID, nil
}

// RemoveTracksFromPlaylist removes one or more tracks from a user's playlist.
// This call requrles that the user has authorized the ScopePlaylistModifyPublic
// or ScopePlaylistModifyPrivate scopes.
//
// If the track(s) occur multiple times in the specified playlist, then all occurrences
// of the track will be removed.  If successful, the snapshot ID returned can be used to
// identify the playlist version in future requests.
func (c *Client) RemoveTracksFromPlaylist(ctx context.Context, playlistID ID, trackIDs ...ID) (newSnapshotID string, err error) {
	tracks := make([]struct {
		URI string `json:"uri"`
	}, len(trackIDs))

	for i, u := range trackIDs {
		tracks[i].URI = fmt.Sprintf("spotify:track:%s", u)
	}
	return c.removeTracksFromPlaylist(ctx, playlistID, tracks, "")
}

// TrackToRemove specifies a track to be removed from a playlist.
// Positions is a slice of 0-based track indices.
// TrackToRemove is used with RemoveTracksFromPlaylistOpt.
type TrackToRemove struct {
	URI       string `json:"uri"`
	Positions []int  `json:"positions"`
}

// NewTrackToRemove creates a new TrackToRemove object with the specified
// track ID and playlist locations.
func NewTrackToRemove(trackID string, positions []int) TrackToRemove {
	return TrackToRemove{
		URI:       fmt.Sprintf("spotify:track:%s", trackID),
		Positions: positions,
	}
}

// RemoveTracksFromPlaylistOpt is like RemoveTracksFromPlaylist, but it supports
// optional parameters that offer more fine-grained control.  Instead of deleting
// all occurrences of a track, this function takes an index with each track URI
// that indicates the position of the track in the playlist.
//
// In addition, the snapshotID parameter allows you to specify the snapshot ID
// against which you want to make the changes.  Spotify will validate that the
// specified tracks exist in the specified positions and make the changes, even
// if more recent changes have been made to the playlist.  If a track in the
// specified position is not found, the entire request will fail and no edits
// will take place. (Note: the snapshot is optional, pass the empty string if
// you don't care about it.)
func (c *Client) RemoveTracksFromPlaylistOpt(
	ctx context.Context,
	playlistID ID,
	tracks []TrackToRemove,
	snapshotID string) (newSnapshotID string, err error) {

	return c.removeTracksFromPlaylist(ctx, playlistID, tracks, snapshotID)
}

func (c *Client) removeTracksFromPlaylist(
	ctx context.Context,
	playlistID ID,
	tracks interface{},
	snapshotID string) (newSnapshotID string, err error) {

	m := make(map[string]interface{})
	m["tracks"] = tracks
	if snapshotID != "" {
		m["snapshot_id"] = snapshotID
	}

	spotifyURL := fmt.Sprintf("%splaylists/%s/tracks",
		c.baseURL, string(playlistID))
	body, err := json.Marshal(m)
	if err != nil {
		return "", err
	}
	req, err := http.NewRequestWithContext(ctx, "DELETE", spotifyURL, bytes.NewReader(body))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")

	result := struct {
		SnapshotID string `json:"snapshot_id"`
	}{}

	err = c.execute(req, &result)
	if err != nil {
		return "", err
	}

	return result.SnapshotID, err
}

// ReplacePlaylistTracks replaces all of the tracks in a playlist, overwriting its
// existing tracks  This can be useful for replacing or reordering tracks, or for
// clearing a playlist.
//
// Modifying a public playlist requires that the user has authorized the
// ScopePlaylistModifyPublic scope.  Modifying a private playlist requires the
// ScopePlaylistModifyPrivate scope.
//
// A maximum of 100 tracks is permited in this call.  Additional tracks must be
// added via AddTracksToPlaylist.
func (c *Client) ReplacePlaylistTracks(ctx context.Context, playlistID ID, trackIDs ...ID) error {
	trackURIs := make([]string, len(trackIDs))
	for i, u := range trackIDs {
		trackURIs[i] = fmt.Sprintf("spotify:track:%s", u)
	}
	spotifyURL := fmt.Sprintf("%splaylists/%s/tracks?uris=%s",
		c.baseURL, playlistID, strings.Join(trackURIs, ","))
	req, err := http.NewRequestWithContext(ctx, "PUT", spotifyURL, nil)
	if err != nil {
		return err
	}
	err = c.execute(req, nil, http.StatusCreated)
	if err != nil {
		return err
	}

	return nil
}

// UserFollowsPlaylist checks if one or more (up to 5) Spotify users are following
// a Spotify playlist, given the playlist's owner and ID.
//
// Checking if a user follows a playlist publicly doesn't require any scopes.
// Checking if the user is privately following a playlist is only possible for the
// current user when that user has granted access to the ScopePlaylistReadPrivate scope.
func (c *Client) UserFollowsPlaylist(ctx context.Context, playlistID ID, userIDs ...string) ([]bool, error) {
	spotifyURL := fmt.Sprintf("%splaylists/%s/followers/contains?ids=%s",
		c.baseURL, playlistID, strings.Join(userIDs, ","))

	follows := make([]bool, len(userIDs))

	err := c.get(ctx, spotifyURL, &follows)
	if err != nil {
		return nil, err
	}

	return follows, err
}

// PlaylistReorderOptions is used with ReorderPlaylistTracks to reorder
// a track or group of tracks in a playlist.
//
// For example, in a playlist with 10 tracks, you can:
//
// - move the first track to the end of the playlist by setting
//   RangeStart to 0 and InsertBefore to 10
// - move the last track to the beginning of the playlist by setting
//   RangeStart to 9 and InsertBefore to 0
// - Move the last 2 tracks to the beginning of the playlist by setting
//   RangeStart to 8 and RangeLength to 2.
type PlaylistReorderOptions struct {
	// The position of the first track to be reordered.
	// This field is required.
	RangeStart int `json:"range_start"`
	// The amount of tracks to be reordered.  This field is optional.  If
	// you don't set it, the value 1 will be used.
	RangeLength int `json:"range_length,omitempty"`
	// The position where the tracks should be inserted.  To reorder the
	// tracks to the end of the playlist, simply set this to the position
	// after the last track.  This field is required.
	InsertBefore int `json:"insert_before"`
	// The playlist's snapshot ID against which you wish to make the changes.
	// This field is optional.
	SnapshotID string `json:"snapshot_id,omitempty"`
}

// ReorderPlaylistTracks reorders a track or group of tracks in a playlist.  It
// returns a snapshot ID that can be used to identify the [newly modified] playlist
// version in future requests.
//
// See the docs for PlaylistReorderOptions for information on how the reordering
// works.
//
// Reordering tracks in the current user's public playlist requires ScopePlaylistModifyPublic.
// Reordering tracks in the user's private playlists (including collaborative playlists) requires
// ScopePlaylistModifyPrivate.
func (c *Client) ReorderPlaylistTracks(ctx context.Context, playlistID ID, opt PlaylistReorderOptions) (snapshotID string, err error) {
	spotifyURL := fmt.Sprintf("%splaylists/%s/tracks", c.baseURL, playlistID)
	j, err := json.Marshal(opt)
	if err != nil {
		return "", err
	}
	req, err := http.NewRequestWithContext(ctx, "PUT", spotifyURL, bytes.NewReader(j))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")

	result := struct {
		SnapshotID string `json:"snapshot_id"`
	}{}
	err = c.execute(req, &result)
	if err != nil {
		return "", err
	}

	return result.SnapshotID, err
}

// SetPlaylistImage replaces the image used to represent a playlist.
// This action can only be performed by the owner of the playlist,
// and requires ScopeImageUpload as well as ScopeModifyPlaylist{Public|Private}..
func (c *Client) SetPlaylistImage(ctx context.Context, playlistID ID, img io.Reader) error {
	spotifyURL := fmt.Sprintf("%splaylists/%s/images", c.baseURL, playlistID)
	// data flow:
	// img (reader) -> copy into base64 encoder (writer) -> pipe (write end)
	// pipe (read end) -> request body
	r, w := io.Pipe()
	go func() {
		enc := base64.NewEncoder(base64.StdEncoding, w)
		_, err := io.Copy(enc, img)
		_ = enc.Close()
		_ = w.CloseWithError(err)
	}()

	req, err := http.NewRequestWithContext(ctx, "PUT", spotifyURL, r)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "image/jpeg")
	return c.execute(req, nil, http.StatusAccepted)
}
