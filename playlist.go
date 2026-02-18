package spotify

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// PlaylistItems contains details about the items in a playlist.
type PlaylistItems struct {
	// A link to the Web API endpoint where full details of
	// the playlist's items can be retrieved.
	Endpoint string `json:"href"`
	// The total number of items in the playlist.
	Total Numeric `json:"total"`
}

// SimplePlaylist contains basic info about a Spotify playlist.
type SimplePlaylist struct {
	// Indicates whether the playlist owner allows others to modify the playlist.
	// Note: only non-collaborative playlists are currently returned by Spotify's Web API.
	Collaborative bool `json:"collaborative"`
	// The playlist description. Empty string if no description is set.
	Description  string            `json:"description"`
	ExternalURLs map[string]string `json:"external_urls"`
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
	// items can be retrieved, along with the total number of items in the playlist.
	Items PlaylistItems `json:"items"`
	URI   URI           `json:"uri"`
}

// FullPlaylist provides extra playlist data in addition to the data provided by [SimplePlaylist].
type FullPlaylist struct {
	SimplePlaylist
	// Information about the followers of this playlist.
	Followers Followers        `json:"followers"`
	Items     PlaylistItemPage `json:"items"`
}

// FeaturedPlaylists gets a [list of playlists featured by Spotify].
//
// Supported options: [Locale], [Country], [Timestamp], [Limit], [Offset].
//
// [list of playlists featured by Spotify]: https://developer.spotify.com/documentation/web-api/reference/get-featured-playlists
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

// GetPlaylist [fetches a playlist] from spotify.
//
// Supported options: [Fields].
//
// [fetches a playlist]: https://developer.spotify.com/documentation/web-api/reference/get-playlist
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

// PlaylistItem contains info about an item in a playlist.
type PlaylistItem struct {
	// The date and time the track was added to the playlist.
	// You can use [TimestampLayout] to convert
	// this field to a [time.Time].
	// Warning: very old playlists may not populate this value.
	AddedAt string `json:"added_at"`
	// The Spotify user who added the track to the playlist.
	// Warning: very old playlists may not populate this value.
	AddedBy User `json:"added_by"`
	// Whether this track is a local file or not.
	IsLocal bool `json:"is_local"`
	// Information about the item.
	Item PlaylistItemTrack `json:"item"`
}

// PlaylistItemTrack is a union type for both tracks and episodes. If both
// values are null, it's likely that the piece of content is not available in
// the configured market.
type PlaylistItemTrack struct {
	Track   *FullTrack
	Episode *EpisodePage
}

// UnmarshalJSON customises the unmarshalling based on the type flags set.
func (t *PlaylistItemTrack) UnmarshalJSON(b []byte) error {
	// Spotify API will return `track: null`` where the content is not available
	// in the specified market. We should respect this and just pass the null
	// up...
	if bytes.Equal(b, []byte("null")) {
		return nil
	}

	itemType := struct {
		Type string `json:"type"`
	}{}

	err := json.Unmarshal(b, &itemType)
	if err != nil {
		return err
	}

	switch itemType.Type {
	case "episode":
		return json.Unmarshal(b, &t.Episode)
	case "track":
		return json.Unmarshal(b, &t.Track)
	default:
		return fmt.Errorf("unrecognized item type: %s", itemType.Type)
	}
}

// PlaylistItemPage contains information about items in a playlist.
type PlaylistItemPage struct {
	basePage
	Items []PlaylistItem `json:"items"`
}

// GetPlaylistItems [gets full details of the items in a playlist], given the
// playlist's [Spotify ID].
//
// Supported options: [Limit], [Offset], [Market], [Fields].
//
// [gets full details of the items in a playlist]: https://developer.spotify.com/documentation/web-api/reference/get-playlists-tracks
// [Spotify ID]: https://developer.spotify.com/documentation/web-api/#spotify-uris-and-ids
func (c *Client) GetPlaylistItems(ctx context.Context, playlistID ID, opts ...RequestOption) (*PlaylistItemPage, error) {
	spotifyURL := fmt.Sprintf("%splaylists/%s/items", c.baseURL, playlistID)

	// Add default as the first option so it gets override by url.Values#Set
	opts = append([]RequestOption{AdditionalTypes(EpisodeAdditionalType, TrackAdditionalType)}, opts...)

	if params := processOptions(opts...).urlParams.Encode(); params != "" {
		spotifyURL += "?" + params
	}

	var result PlaylistItemPage

	err := c.get(ctx, spotifyURL, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

// CreatePlaylist [creates a playlist] for the current user.
// The playlist will be empty until you add items to it.
// The playlistName does not need to be unique - a user can have
// several playlists with the same name.
//
// Creating a public playlist requires [ScopePlaylistModifyPublic];
// creating a private playlist requires [ScopePlaylistModifyPrivate].
//
// On success, the newly created playlist is returned.
//
// [creates a playlist]: https://developer.spotify.com/documentation/web-api/reference/create-playlist
func (c *Client) CreatePlaylist(ctx context.Context, playlistName, description string, public bool, collaborative bool) (*FullPlaylist, error) {
	spotifyURL := fmt.Sprintf("%sme/playlists", c.baseURL)
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

	return &p, nil
}

// ChangePlaylistName [changes the name of a playlist].  This call requires that the
// user has authorized the [ScopePlaylistModifyPublic] or [ScopePlaylistModifyPrivate]
// scopes (depending on whether the playlist is public or private).
// The current user must own the playlist in order to modify it.
//
// [changes the name of a playlist]: https://developer.spotify.com/documentation/web-api/reference/change-playlist-details
func (c *Client) ChangePlaylistName(ctx context.Context, playlistID ID, newName string) error {
	return c.modifyPlaylist(ctx, playlistID, newName, "", nil)
}

// ChangePlaylistAccess [modifies the public/private status of a playlist].  This call
// requires that the user has authorized the [ScopePlaylistModifyPublic] or
// [ScopePlaylistModifyPrivate] scopes (depending on whether the playlist is
// currently public or private).  The current user must own the playlist to modify it.
//
// [modifies the public/private status of a playlist]: https://developer.spotify.com/documentation/web-api/reference/change-playlist-details
func (c *Client) ChangePlaylistAccess(ctx context.Context, playlistID ID, public bool) error {
	return c.modifyPlaylist(ctx, playlistID, "", "", &public)
}

// ChangePlaylistDescription [modifies the description of a playlist].  This call
// requires that the user has authorized the [ScopePlaylistModifyPublic] or
// [ScopePlaylistModifyPrivate] scopes (depending on whether the playlist is
// currently public or private).  The current user must own the playlist to modify it.
//
// [modifies the description of a playlist]: https://developer.spotify.com/documentation/web-api/reference/change-playlist-details
func (c *Client) ChangePlaylistDescription(ctx context.Context, playlistID ID, newDescription string) error {
	return c.modifyPlaylist(ctx, playlistID, "", newDescription, nil)
}

// ChangePlaylistNameAndAccess combines [ChangePlaylistName] and [ChangePlaylistAccess] into
// a single Web API call.  It requires that the user has authorized the [ScopePlaylistModifyPublic]
// or [ScopePlaylistModifyPrivate] scopes (depending on whether the playlist is currently
// public or private).  The current user must own the playlist to modify it.
func (c *Client) ChangePlaylistNameAndAccess(ctx context.Context, playlistID ID, newName string, public bool) error {
	return c.modifyPlaylist(ctx, playlistID, newName, "", &public)
}

// ChangePlaylistNameAccessAndDescription combines [ChangePlaylistName], [ChangePlaylistAccess], and
// [ChangePlaylistDescription] into a single Web API call.  It requires that the user has authorized
// the [ScopePlaylistModifyPublic] or [ScopePlaylistModifyPrivate] scopes (depending on whether the
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
	return c.execute(req, nil, http.StatusCreated)
}

// AddTracksToPlaylist [adds one or more tracks to a user's playlist].
// This call requires [ScopePlaylistModifyPublic] or [ScopePlaylistModifyPrivate].
// A maximum of 100 tracks can be added per call.  It returns a snapshot ID that
// can be used to identify this version (the new version) of the playlist in
// future requests.
//
// [adds one or more tracks to a user's playlist]: https://developer.spotify.com/documentation/web-api/reference/add-tracks-to-playlist
func (c *Client) AddTracksToPlaylist(ctx context.Context, playlistID ID, trackIDs ...ID) (snapshotID string, err error) {
	uris := make([]string, len(trackIDs))
	for i, id := range trackIDs {
		uris[i] = fmt.Sprintf("spotify:track:%s", id)
	}
	m := make(map[string]interface{})
	m["uris"] = uris

	spotifyURL := fmt.Sprintf("%splaylists/%s/items",
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

// RemoveTracksFromPlaylist [removes one or more tracks from a user's playlist].
// This call requires that the user has authorized the [ScopePlaylistModifyPublic]
// or [ScopePlaylistModifyPrivate] scopes.
//
// If the track(s) occur multiple times in the specified playlist, then all occurrences
// of the track will be removed.  If successful, the snapshot ID returned can be used to
// identify the playlist version in future requests.
//
// [removes one or more tracks from a user's playlist]: https://developer.spotify.com/documentation/web-api/reference/remove-tracks-playlist
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

// NewTrackToRemove returns a [TrackToRemove] with the specified
// track ID and playlist locations.
func NewTrackToRemove(trackID string, positions []int) TrackToRemove {
	return TrackToRemove{
		URI:       fmt.Sprintf("spotify:track:%s", trackID),
		Positions: positions,
	}
}

// RemoveTracksFromPlaylistOpt is like [RemoveTracksFromPlaylist], but it supports
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
	snapshotID string,
) (newSnapshotID string, err error) {
	return c.removeTracksFromPlaylist(ctx, playlistID, tracks, snapshotID)
}

func (c *Client) removeTracksFromPlaylist(
	ctx context.Context,
	playlistID ID,
	tracks interface{},
	snapshotID string,
) (newSnapshotID string, err error) {
	m := make(map[string]interface{})
	m["items"] = tracks
	if snapshotID != "" {
		m["snapshot_id"] = snapshotID
	}

	spotifyURL := fmt.Sprintf("%splaylists/%s/items",
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

// ReplacePlaylistTracks [replaces all of the tracks in a playlist], overwriting its
// existing tracks  This can be useful for replacing or reordering tracks, or for
// clearing a playlist.
//
// Modifying a public playlist requires that the user has authorized the
// [ScopePlaylistModifyPublic] scope.  Modifying a private playlist requires the
// [ScopePlaylistModifyPrivate] scope.
//
// A maximum of 100 tracks are permitted in this call.  Additional tracks must be
// added via [AddTracksToPlaylist].
//
// [replaces all of the tracks in a playlist]: https://developer.spotify.com/documentation/web-api/reference/reorder-or-replace-playlists-tracks
func (c *Client) ReplacePlaylistTracks(ctx context.Context, playlistID ID, trackIDs ...ID) error {
	trackURIs := make([]string, len(trackIDs))
	for i, u := range trackIDs {
		trackURIs[i] = fmt.Sprintf("spotify:track:%s", u)
	}
	spotifyURL := fmt.Sprintf("%splaylists/%s/items?uris=%s",
		c.baseURL, playlistID, strings.Join(trackURIs, ","))
	req, err := http.NewRequestWithContext(ctx, "PUT", spotifyURL, nil)
	if err != nil {
		return err
	}
	return c.execute(req, nil, http.StatusCreated)
}

// ReplacePlaylistItems [replaces all the items in a playlist], overwriting its
// existing tracks  This can be useful for replacing or reordering tracks, or for
// clearing a playlist.
//
// Modifying a public playlist requires that the user has authorized the
// [ScopePlaylistModifyPublic] scope.  Modifying a private playlist requires the
// [ScopePlaylistModifyPrivate] scope.
//
// A maximum of 100 tracks is permited in this call.  Additional tracks must be
// added via AddTracksToPlaylist.
//
// [replaces all the items in a playlist]: https://developer.spotify.com/documentation/web-api/reference/reorder-or-replace-playlists-tracks
func (c *Client) ReplacePlaylistItems(ctx context.Context, playlistID ID, items ...URI) (string, error) {
	m := make(map[string]interface{})
	m["uris"] = items

	body, err := json.Marshal(m)
	if err != nil {
		return "", err
	}

	spotifyURL := fmt.Sprintf("%splaylists/%s/items", c.baseURL, playlistID)
	req, err := http.NewRequestWithContext(ctx, "PUT", spotifyURL, bytes.NewReader(body))
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

// PlaylistReorderOptions is used with ReorderPlaylistTracks to reorder
// a track or group of tracks in a playlist.
//
// For example, in a playlist with 10 tracks, you can:
//
//   - move the first track to the end of the playlist by setting
//     RangeStart to 0 and InsertBefore to 10
//   - move the last track to the beginning of the playlist by setting
//     RangeStart to 9 and InsertBefore to 0
//   - Move the last 2 tracks to the beginning of the playlist by setting
//     RangeStart to 8 and RangeLength to 2.
type PlaylistReorderOptions struct {
	// The position of the first track to be reordered.
	// This field is required.
	RangeStart Numeric `json:"range_start"`
	// The amount of tracks to be reordered.  This field is optional.  If
	// you don't set it, the value 1 will be used.
	RangeLength Numeric `json:"range_length,omitempty"`
	// The position where the tracks should be inserted.  To reorder the
	// tracks to the end of the playlist, simply set this to the position
	// after the last track.  This field is required.
	InsertBefore Numeric `json:"insert_before"`
	// The playlist's snapshot ID against which you wish to make the changes.
	// This field is optional.
	SnapshotID string `json:"snapshot_id,omitempty"`
}

// ReorderPlaylistTracks reorders a track or group of tracks in a playlist.  It
// returns a snapshot ID that can be used to identify the (newly modified) playlist
// version in future requests.
//
// See the docs for [PlaylistReorderOptions] for information on how the reordering
// works.
//
// Reordering tracks in the current user's public playlist requires [ScopePlaylistModifyPublic].
// Reordering tracks in the user's private playlists (including collaborative playlists) requires
// [ScopePlaylistModifyPrivate].
func (c *Client) ReorderPlaylistTracks(ctx context.Context, playlistID ID, opt PlaylistReorderOptions) (snapshotID string, err error) {
	spotifyURL := fmt.Sprintf("%splaylists/%s/items", c.baseURL, playlistID)
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

	return result.SnapshotID, nil
}

// SetPlaylistImage replaces the image used to represent a playlist.
// This action can only be performed by the owner of the playlist,
// and requires [ScopeImageUpload] as well as [ScopeModifyPlaylistPublic] or
// [ScopeModifyPlaylistPrivate].
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
