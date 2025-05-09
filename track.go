package spotify

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"
)

type TrackExternalIDs struct {
	ISRC string `json:"isrc"`
	EAN  string `json:"ean"`
	UPC  string `json:"upc"`
}

// SimpleTrack contains basic info about a track.
type SimpleTrack struct {
	Album   SimpleAlbum    `json:"album"`
	Artists []SimpleArtist `json:"artists"`
	// A list of the countries in which the track can be played,
	// identified by their [ISO 3166-1 alpha-2] codes.
	//
	// [ISO 3166-1 alpha=2]: https://en.wikipedia.org/wiki/ISO_3166-1_alpha-2
	AvailableMarkets []string `json:"available_markets"`
	// The disc number (usually 1 unless the album consists of more than one disc).
	DiscNumber Numeric `json:"disc_number"`
	// The length of the track, in milliseconds.
	Duration Numeric `json:"duration_ms"`
	// Whether or not the track has explicit lyrics.
	// true => yes, it does; false => no, it does not.
	Explicit bool `json:"explicit"`
	// External URLs for this track.
	ExternalURLs map[string]string `json:"external_urls"`
	// ExternalIDs are IDs for this track in other databases
	ExternalIDs TrackExternalIDs `json:"external_ids"`
	// A link to the Web API endpoint providing full details for this track.
	Endpoint string `json:"href"`
	ID       ID     `json:"id"`
	Name     string `json:"name"`
	// A URL to a 30 second preview (MP3) of the track.
	PreviewURL string `json:"preview_url"`
	// The number of the track.  If an album has several
	// discs, the track number is the number on the specified
	// DiscNumber.
	TrackNumber Numeric `json:"track_number"`
	URI         URI     `json:"uri"`
	// Type of the track
	Type string `json:"type"`
}

func (st SimpleTrack) String() string {
	return fmt.Sprintf("TRACK<[%s] [%s]>", st.ID, st.Name)
}

// LinkedFromInfo is included in a track response when [Track Relinking] is applied.
//
// [Track Relinking]: https://developer.spotify.com/documentation/general/guides/track-relinking-guide/
type LinkedFromInfo struct {
	// ExternalURLs are the known external APIs for this track or album
	ExternalURLs map[string]string `json:"external_urls"`

	// Href is a link to the Web API endpoint providing full details
	Href string `json:"href"`

	// ID of the linked track
	ID ID `json:"id"`

	// Type of the link: album of the track
	Type string `json:"type"`

	// URI is the [Spotify URI] of the track/album.
	//
	// [Spotify URI]: https://developer.spotify.com/documentation/web-api/#spotify-uris-and-ids
	URI string `json:"uri"`
}

// FullTrack provides extra track data in addition to what is provided by [SimpleTrack].
type FullTrack struct {
	SimpleTrack
	// Popularity of the track.  The value will be between 0 and 100,
	// with 100 being the most popular.  The popularity is calculated from
	// both total plays and most recent plays.
	Popularity Numeric `json:"popularity"`

	// IsPlayable is included when [Track Relinking] is applied, and reports if
	// the track is playable. It's reported when the "market" parameter is
	// passed to the tracks listing API.
	//
	// [Track Relinking]: https://developer.spotify.com/documentation/general/guides/track-relinking-guide/
	IsPlayable *bool `json:"is_playable"`

	// LinkedFromInfo is included in a track response when [Track Relinking] is
	// applied, and points to the linked track. It's reported when the "market"
	// parameter is passed to the tracks listing API.
	//
	// [Track Relinking]: https://developer.spotify.com/documentation/general/guides/track-relinking-guide/
	LinkedFrom *LinkedFromInfo `json:"linked_from"`
}

// PlaylistTrack contains info about a track in a playlist.
type PlaylistTrack struct {
	// The date and time the track was added to the playlist. You can use
	// [TimestampLayout] to convert this field to a [time.Time].
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
	// The date and time the track was saved, represented as an ISO 8601 UTC
	// timestamp with a zero offset (YYYY-MM-DDTHH:MM:SSZ). You can use
	// [TimestampLayout] to convert this to a [time.Time].
	AddedAt   string `json:"added_at"`
	FullTrack `json:"track"`
}

// TimeDuration returns the track's duration as a [time.Duration] value.
func (t *SimpleTrack) TimeDuration() time.Duration {
	return time.Duration(t.Duration) * time.Millisecond
}

// GetTrack gets Spotify catalog information for
// a [single track] identified by its unique [Spotify ID].
//
// Supported options: [Market].
//
// [single track]: https://developer.spotify.com/documentation/web-api/reference/get-track
// [Spotify ID]: https://developer.spotify.com/documentation/web-api/#spotify-uris-and-ids
func (c *Client) GetTrack(ctx context.Context, id ID, opts ...RequestOption) (*FullTrack, error) {
	spotifyURL := c.baseURL + "tracks/" + string(id)

	var t FullTrack

	if params := processOptions(opts...).urlParams.Encode(); params != "" {
		spotifyURL += "?" + params
	}

	err := c.get(ctx, spotifyURL, &t)
	if err != nil {
		return nil, err
	}

	return &t, nil
}

// GetTracks gets Spotify catalog information for [multiple tracks] based on their
// Spotify IDs.  It supports up to 50 tracks in a single call.  Tracks are
// returned in the order requested.  If a track is not found, that position in the
// result will be nil.  Duplicate ids in the query will result in duplicate
// tracks in the result.
//
// Supported options: [Market].
//
// [multiple tracks]: https://developer.spotify.com/documentation/web-api/reference/get-several-tracks
func (c *Client) GetTracks(ctx context.Context, ids []ID, opts ...RequestOption) ([]*FullTrack, error) {
	if len(ids) > 50 {
		return nil, errors.New("spotify: FindTracks supports up to 50 tracks")
	}

	params := processOptions(opts...).urlParams
	params.Set("ids", strings.Join(toStringSlice(ids), ","))
	spotifyURL := c.baseURL + "tracks?" + params.Encode()

	var t struct {
		Tracks []*FullTrack `json:"tracks"`
	}

	err := c.get(ctx, spotifyURL, &t)
	if err != nil {
		return nil, err
	}

	return t.Tracks, nil
}
