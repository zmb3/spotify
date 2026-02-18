package spotify

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"
)

// SimpleAlbum contains basic data about an album.
type SimpleAlbum struct {
	// The name of the album.
	Name string `json:"name"`
	// A slice of [SimpleArtist].
	Artists []SimpleArtist `json:"artists"`

	// The type of the album: one of "album",
	// "single", or "compilation".
	AlbumType string `json:"album_type"`
	// The [Spotify ID] for the album.
	//
	// [Spotify ID]: https://developer.spotify.com/documentation/web-api/concepts/spotify-uris-ids
	ID ID `json:"id"`
	// The [Spotify URI] for the album.
	//
	// [Spotify URI]: https://developer.spotify.com/documentation/web-api/concepts/spotify-uris-ids
	URI URI `json:"uri"`
	// A link to the Web API endpoint providing full
	// details of the album.
	Endpoint string `json:"href"`
	// The cover art for the album in various sizes,
	// widest first.
	Images []Image `json:"images"`
	// Known external URLs for this album.
	ExternalURLs map[string]string `json:"external_urls"`
	// The date the album was first released.  For example, "1981-12-15".
	// Depending on the ReleaseDatePrecision, it might be shown as
	// "1981" or "1981-12". You can use [SimpleAlbum.ReleaseDateTime] to convert
	// this to a [time.Time] value.
	ReleaseDate string `json:"release_date"`
	// The precision with which ReleaseDate value is known: "year", "month", or "day"
	ReleaseDatePrecision string `json:"release_date_precision"`
	// The number of tracks on the album.
	TotalTracks Numeric `json:"total_tracks"`
}

// ReleaseDateTime converts [SimpleAlbum.ReleaseDate] to a [time.Time].
// All of the fields in the result may not be valid.  For example, if
// [SimpleAlbum.ReleaseDatePrecision] is "month", then only the month and year
// (but not the day) of the result are valid.
func (s *SimpleAlbum) ReleaseDateTime() time.Time {
	if s.ReleaseDatePrecision == "day" {
		result, _ := time.Parse(DateLayout, s.ReleaseDate)
		return result
	}
	if s.ReleaseDatePrecision == "month" {
		ym := strings.Split(s.ReleaseDate, "-")
		year, _ := strconv.Atoi(ym[0])
		month, _ := strconv.Atoi(ym[1])
		return time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
	}
	year, _ := strconv.Atoi(s.ReleaseDate)
	return time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC)
}

// Copyright contains the copyright statement associated with an album.
type Copyright struct {
	// The copyright text for the album.
	Text string `json:"text"`
	// The type of copyright.
	Type string `json:"type"`
}

// FullAlbum provides extra album data in addition to the data provided by [SimpleAlbum].
type FullAlbum struct {
	SimpleAlbum
	Copyrights []Copyright     `json:"copyrights"`
	Genres     []string        `json:"genres"`
	Tracks     SimpleTrackPage `json:"tracks"`
}

// SavedAlbum provides info about an album saved to a user's account.
type SavedAlbum struct {
	// The date and time the track was saved, represented as an ISO
	// 8601 UTC timestamp with a zero offset (YYYY-MM-DDTHH:MM:SSZ).
	// You can use [TimestampLayout] to convert this to a [time.Time] value.
	AddedAt   string `json:"added_at"`
	FullAlbum `json:"album"`
}

// GetAlbum gets Spotify catalog information for a single album, given its
// [Spotify ID]. Supported options: [Market].
//
// [Spotify ID]: https://developer.spotify.com/documentation/web-api/concepts/spotify-uris-ids
func (c *Client) GetAlbum(ctx context.Context, id ID, opts ...RequestOption) (*FullAlbum, error) {
	spotifyURL := fmt.Sprintf("%salbums/%s", c.baseURL, id)

	if params := processOptions(opts...).urlParams.Encode(); params != "" {
		spotifyURL += "?" + params
	}

	var a FullAlbum

	err := c.get(ctx, spotifyURL, &a)
	if err != nil {
		return nil, err
	}

	return &a, nil
}

func toStringSlice(ids []ID) []string {
	result := make([]string, len(ids))
	for i, str := range ids {
		result[i] = str.String()
	}
	return result
}

// AlbumType represents the type of an album. It can be used to filter
// results when searching for albums.
type AlbumType int

// AlbumType values that can be used to filter which types of albums are
// searched for.  These are flags that can be bitwise OR'd together
// to search for multiple types of albums simultaneously.
const (
	AlbumTypeAlbum AlbumType = 1 << iota
	AlbumTypeSingle
	AlbumTypeAppearsOn
	AlbumTypeCompilation
)

func (at AlbumType) encode() string {
	types := []string{}
	if at&AlbumTypeAlbum != 0 {
		types = append(types, "album")
	}
	if at&AlbumTypeSingle != 0 {
		types = append(types, "single")
	}
	if at&AlbumTypeAppearsOn != 0 {
		types = append(types, "appears_on")
	}
	if at&AlbumTypeCompilation != 0 {
		types = append(types, "compilation")
	}
	return strings.Join(types, ",")
}

// GetAlbumTracks gets the [tracks] for a particular album.
// If you only care about the tracks, this call is more efficient
// than [GetAlbum].
//
// Supported Options: [Market], [Limit], [Offset].
//
// [tracks]: https://developer.spotify.com/documentation/web-api/reference/get-an-albums-tracks
func (c *Client) GetAlbumTracks(ctx context.Context, id ID, opts ...RequestOption) (*SimpleTrackPage, error) {
	spotifyURL := fmt.Sprintf("%salbums/%s/tracks", c.baseURL, id)

	if params := processOptions(opts...).urlParams.Encode(); params != "" {
		spotifyURL += "?" + params
	}

	var result SimpleTrackPage
	err := c.get(ctx, spotifyURL, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}
