package spotify

// PlaylistTracks contains details about the tracks in a
// playlist.
type PlaylistTracks struct {
	// A link to the Web API endpoint where full details of
	// the playlist's tracks can be retrieved.
	Endpoint string `json:"href"`
	// The total number of tracks in the playlist.
	Total uint `json:"total"`
}

// SimplePlaylist contains basic info about a
// Spotify playlist.
type SimplePlaylist struct {
	// Indicates whether the playlist owner allows
	// others to modify the playlist.  Note: only
	// non-collaborative playlists are currently
	// returned by Spotify's Web API.
	Collaborative bool `json:"collaborative"`
	// Known external URLs for this playlist.
	ExternalURLs ExternalURL `json:"external_urls"`
	// A link to the Web API endpoint providing full
	// details of the playlist.
	Endpoint string `json:"href"`
	// The Spotify ID for the playlist.
	ID ID `json:"id"`
	// The playlist image.  Note: this field is only
	// returned for modified, verified playlists.
	// Otherwise the slice is empty.  If returned,
	// the source URL for the image is temporary
	// and will expire in less than a day.
	Images []Image `json:"images"`
	// The name of the playlist.
	Name string `json:"name"`
	// The user who owns the playlist.
	Owner User `json:"owner"`
	// The playlist's public/private status
	IsPublic bool `json:"public"`
	// A collection to the Web API endpoint where
	// full details of the playlist's tracks can be
	// retrieved, along with the total number of
	// tracks in the playlist.
	Tracks PlaylistTracks `json:"tracks"`
	// The Spotify URI for the playlist.
	URI URI `json:"uri"`
}

// FullPlaylist provides extra playlist data in addition
// to the data provided by SimplePlaylist.
type FullPlaylist struct {
	SimplePlaylist
	// The playlist description.  Only returned for modified,
	// verified playlists.
	Description string `json:"description"`
	// Information about the followers of this playlist.
	Followers Followers `json:"followers"`
	// The version identifier for the current playlist.
	// Can be supplied in other requests to target a
	// specific playlist version.
	SnapshotID string `json:"snapshot_id"`
	// Information about the tracks of the playlist.
	// TODO: array of playlist track objects inside a
	// TODO: paging object.  is this the same as simple?
	Tracks string `json:"tracks"`
}

func (p *SimplePlaylist) String() string {
	return "SimplePlaylist: " + p.Name
}

func (p *FullPlaylist) String() string {
	return "FullPlaylist: " + p.Name
}
