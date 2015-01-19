package spotify

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

// SimpleAlbum contains basic data about an album.
type SimpleAlbum struct {
	// The name of the album.
	Name string `json:"name"`
	// The type of the album: one of "album",
	// "single", or "compilation".
	AlbumType string `json:"album_type"`
	// The SpotifyID for the album.
	ID ID `json:"id"`
	// The SpotifyURI for the album.
	URI URI `json:"uri"`
	// The markets in which the album is available,
	// identified using ISO 3166-1 alpha-2 country
	// codes.  Note that al album is considered
	// available in a market when at least 1 of its
	// tracks is available in that market.
	AvailableMarkets []string `json:"available_markets"`
	// A link to the Web API enpoint providing full
	// details of the album.
	Endpoint string `json:"href"`
	// The cover art for the album in various sizes,
	// widest first.
	Images []Image `json:"images"`
	// Known external URLs for this album.
	ExternalURLs ExternalURL `json:"external_urls"`
}

// Copyright contains the copyright statement
// associated with an album.
type Copyright struct {
	// The copyright text for the album.
	Text string `json:"text"`
	// The type of copyright.
	Type string `json:"type"`
}

// FullAlbum provides extra album data in addition
// to the data provided by SimpleAlbum.
type FullAlbum struct {
	SimpleAlbum
	// The artists of the album.
	Artists []SimpleArtist `json:"artists"`
	// The copyright statements of the album.
	Copyrights []Copyright `json:"copyrights"`
	// A list of genres used to classify the album.
	// For example, "Prog Rock" or "Post-Grunge".
	// If not yet classified, the slice is empty.
	Genres []string `json:"genres"`
	// The popularity of the album.  The value will
	// be between 0 and 100, with 100 being the most
	// popular.  Popularity of an album is calculated
	// from the popularify of the album's individual
	// tracks.
	Popularity int `json:"popularity"`
	// The date the album was first released.  For
	// example, "1981-12-15".  Depending on the
	// ReleaseDatePrecision, it might be shown as
	// "1981" or "1981-12".
	ReleaseDate string `json:"release_date"` // TODO change to Timestamp
	// The precision with which ReleaseDate value
	// is known: "year", "month", or "day"
	ReleaseDatePrecision string `json:"release_date_precision"`
	// The tracks of the album.  Tracks are inside a paging object.
	Tracks TrackResult `json:"tracks"`
	// Known external IDs for the album.
	ExternalIDs ExternalID `json:"external_ids"`
}

// FindAlbum gets Spotify catalog information for a single
// album, given that album's Spotify ID.
func (c *Client) FindAlbum(id ID) (*FullAlbum, error) {
	uri := baseAddress + "albums/" + string(id)
	resp, err := c.http.Get(string(uri))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, decodeError(resp.Body)
	}
	var a FullAlbum
	err = json.NewDecoder(resp.Body).Decode(&a)
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

// FindAlbums is a wrapper around DefaultClient.FindAlbums.
func FindAlbums(ids ...ID) ([]*FullAlbum, error) {
	return DefaultClient.FindAlbums(ids...)
}

// FindAlbums gets Spotify Catalog information for multiple
// albums, given their Spotify IDs.  It supports up to 20
// IDs in a single call.  Albums are returned in the order
// requested.  If an album is not found, that position in the
// result slice will be nil.
func (c *Client) FindAlbums(ids ...ID) ([]*FullAlbum, error) {
	if len(ids) > 20 {
		return nil, errors.New("spotify: exceeded maximum number of albums")
	}
	uri := baseAddress + "albums?ids=" + strings.Join(toStringSlice(ids), ",")
	resp, err := c.http.Get(uri)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, decodeError(resp.Body)
	}
	var a struct {
		Albums []*FullAlbum `json:"albums"`
	}
	err = json.NewDecoder(resp.Body).Decode(&a)
	if err != nil {
		return nil, err
	}
	return a.Albums, nil
}

// AlbumType represents the type of an album.
// It can be used to filter results when searching
// for albums.
type AlbumType int

// AlbumType values that can be used to filter
// which types of albums are searched for.
// These are flags that can be bitwise OR'd together
//  to search for multiple types of albums simultaneously.
const (
	AlbumTypeAlbum       AlbumType = 1 << iota
	AlbumTypeSingle                = 1 << iota
	AlbummTypeAppearsOn            = 1 << iota
	AlbumTypeCompilation           = 1 << iota
)

func (at AlbumType) encode() string {
	types := []string{}
	if at&AlbumTypeAlbum != 0 {
		types = append(types, "album")
	}
	if at&AlbumTypeSingle != 0 {
		types = append(types, "single")
	}
	if at&AlbummTypeAppearsOn != 0 {
		types = append(types, "appears_on")
	}
	if at&AlbumTypeCompilation != 0 {
		types = append(types, "compilation")
	}
	return strings.Join(types, ",")
}

// FindAlbumTracks is a wrapper around DefaultClient.FindAlbumTracks.
func FindAlbumTracks(id ID) (*TrackResult, error) {
	return DefaultClient.FindAlbumTracks(id)
}

// FindAlbumTracks gets the tracks for a particular album.
// If you only care about the tracks, this call is more efficient
// than FindAlbum.
func (c *Client) FindAlbumTracks(id ID) (*TrackResult, error) {
	return c.FindAlbumTracksLimited(id, -1, -1)
}

// FindAlbumTracksLimited is a wrapper around DefaultClient.FindAlbumTracksLimited.
func FindAlbumTracksLimited(id ID, limit, offset int) (*TrackResult, error) {
	return DefaultClient.FindAlbumTracksLimited(id, limit, offset)
}

// FindAlbumTracksLimited behaves like FindAlbumTracks, with the
// exception that it allows you to specify extra parameters that
// limit the number of results returned.
// The maximum number of results to return is specified by limit.
// The offset argument can be used to specify the index of the first
// track to return.  It can be used along with limit to reqeust
// the next set of results.
func (c *Client) FindAlbumTracksLimited(id ID, limit, offset int) (*TrackResult, error) {
	uri := baseAddress + "albums/" + string(id) + "/tracks"
	v := url.Values{}
	if limit != -1 {
		v.Set("limit", strconv.Itoa(limit))
	}
	if offset != -1 {
		v.Set("offset", strconv.Itoa(offset))
	}
	optional := v.Encode()
	if optional != "" {
		uri = uri + "?" + optional
	}
	resp, err := c.http.Get(uri)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var p page
	err = json.NewDecoder(resp.Body).Decode(&p)
	if err != nil {
		return nil, err
	}
	var result TrackResult
	result.FullResult = p.Endpoint
	result.Limit = p.Limit
	result.Offset = p.Offset
	result.Next = p.Next
	result.Total = p.Total
	result.Previous = p.Previous

	err = json.Unmarshal([]byte(p.Items), &result.Tracks)
	if err != nil {
		return nil, err
	}
	return &result, nil
}
