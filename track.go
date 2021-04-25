package spotify

import (
	"errors"
	"fmt"
	"net/url"
	"strings"
	"time"
)

// SimpleTrack contains basic info about a track.
type SimpleTrack struct {
	Artists []SimpleArtist `json:"artists"`
	// A list of the countries in which the track can be played,
	// identified by their ISO 3166-1 alpha-2 codes.
	AvailableMarkets []string `json:"available_markets"`
	// The disc number (usually 1 unless the album consists of more than one disc).
	DiscNumber int `json:"disc_number"`
	// The length of the track, in milliseconds.
	Duration int `json:"duration_ms"`
	// Whether or not the track has explicit lyrics.
	// true => yes, it does; false => no, it does not.
	Explicit bool `json:"explicit"`
	// External URLs for this track.
	ExternalURLs map[string]string `json:"external_urls"`
	// A link to the Web API endpoint providing full details for this track.
	Endpoint string `json:"href"`
	ID       ID     `json:"id"`
	Name     string `json:"name"`
	// A URL to a 30 second preview (MP3) of the track.
	PreviewURL string `json:"preview_url"`
	// The number of the track.  If an album has several
	// discs, the track number is the number on the specified
	// DiscNumber.
	TrackNumber int `json:"track_number"`
	URI         URI `json:"uri"`
	// Type of the track
	Type string `json:"type"`
}

func (st SimpleTrack) String() string {
	return fmt.Sprintf("TRACK<[%s] [%s]>", st.ID, st.Name)
}

// LinkedFromInfo
// See: https://developer.spotify.com/documentation/general/guides/track-relinking-guide/
type LinkedFromInfo struct {
	// ExternalURLs are the known external APIs for this track or album
	ExternalURLs map[string]string `json:"external_urls"`

	// Href is a link to the Web API endpoint providing full details
	Href string `json:"href"`

	// ID of the linked track
	ID ID `json:"id"`

	// Type of the link: album of the track
	Type string `json:"type"`

	// URI is the Spotify URI of the track/album
	URI string `json:"uri"`
}

// FullTrack provides extra track data in addition to what is provided by SimpleTrack.
type FullTrack struct {
	SimpleTrack
	// The album on which the track appears. The album object includes a link in href to full information about the album.
	Album SimpleAlbum `json:"album"`
	// Known external IDs for the track.
	ExternalIDs map[string]string `json:"external_ids"`
	// Popularity of the track.  The value will be between 0 and 100,
	// with 100 being the most popular.  The popularity is calculated from
	// both total plays and most recent plays.
	Popularity int `json:"popularity"`

	// IsPlayable defines if the track is playable. It's reported when the "market" parameter is passed to the tracks
	// listing API.
	// See: https://developer.spotify.com/documentation/general/guides/track-relinking-guide/
	IsPlayable *bool `json:"is_playable"`

	// LinkedFrom points to the linked track. It's reported when the "market" parameter is passed to the tracks listing
	// API.
	LinkedFrom *LinkedFromInfo `json:"linked_from"`
}

// PlaylistTrack contains info about a track in a playlist.
type PlaylistTrack struct {
	// The date and time the track was added to the playlist.
	// You can use the TimestampLayout constant to convert
	// this field to a time.Time value.
	// Warning: very old playlists may not populate this value.
	AddedAt string `json:"added_at"`
	// The Spotify user who added the track to the playlist.
	// Warning: vary old playlists may not populate this value.
	AddedBy User `json:"added_by"`
	// Whether this track is a local file or not.
	IsLocal bool `json:"is_local"`
	// Information about the track.
	Track FullTrack `json:"track"`
}

// SavedTrack provides info about a track saved to a user's account.
type SavedTrack struct {
	// The date and time the track was saved, represented as an ISO
	// 8601 UTC timestamp with a zero offset (YYYY-MM-DDTHH:MM:SSZ).
	// You can use the TimestampLayout constant to convert this to
	// a time.Time value.
	AddedAt   string `json:"added_at"`
	FullTrack `json:"track"`
}

// TimeDuration returns the track's duration as a time.Duration value.
func (t *SimpleTrack) TimeDuration() time.Duration {
	return time.Duration(t.Duration) * time.Millisecond
}

// GetTrack gets Spotify catalog information for
// a single track identified by its unique Spotify ID.
// API Doc: https://developer.spotify.com/documentation/web-api/reference/tracks/get-track/
func (c *Client) GetTrack(id ID) (*FullTrack, error) {
	return c.GetTrackOpt(id, nil)
}

// GetTrackOpt is like GetTrack but it accepts additional arguments
func (c *Client) GetTrackOpt(id ID, opt *Options) (*FullTrack, error) {
	spotifyURL := c.baseURL + "tracks/" + string(id)

	var t FullTrack

	if opt != nil {
		v := url.Values{}
		if opt.Country != nil {
			v.Set("market", *opt.Country)
		}
		if params := v.Encode(); params != "" {
			spotifyURL += "?" + params
		}
	}

	err := c.get(spotifyURL, &t)
	if err != nil {
		return nil, err
	}

	return &t, nil
}

// GetTracks gets Spotify catalog information for multiple tracks based on their
// Spotify IDs.  It supports up to 50 tracks in a single call.  Tracks are
// returned in the order requested.  If a track is not found, that position in the
// result will be nil.  Duplicate ids in the query will result in duplicate
// tracks in the result.
// API Doc: https://developer.spotify.com/documentation/web-api/reference/tracks/get-several-tracks/
func (c *Client) GetTracks(ids ...ID) ([]*FullTrack, error) {
	return c.GetTracksOpt(nil, ids...)
}

// GetTracksOpt is like GetTracks but it accepts an additional country option for track relinking
func (c *Client) GetTracksOpt(opt *Options, ids ...ID) ([]*FullTrack, error) {
	if len(ids) > 50 {
		return nil, errors.New("spotify: FindTracks supports up to 50 tracks")
	}

	params := url.Values{}
	params.Set("ids", strings.Join(toStringSlice(ids), ","))
	if opt != nil && opt.Country != nil {
		params.Set("market", *opt.Country)
	}
	spotifyURL := c.baseURL + "tracks?" + params.Encode()

	var t struct {
		Tracks []*FullTrack `json:"tracks"`
	}

	err := c.get(spotifyURL, &t)
	if err != nil {
		return nil, err
	}

	return t.Tracks, nil
}
