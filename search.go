package spotify

import (
	"encoding/json"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

const (
	// MarketFromToken can be used in place of the Market parameter
	// if the Client has a valid access token.  In this case, the
	// results will be limited to content that is playable in the
	// country associated with the user's account.  The user must have
	// granted access to the user-read-private scope when the access
	// token was issued.
	MarketFromToken = "from_token"
)

// SearchType represents the type of a query used
// in the Search function.
type SearchType int

// Search type values that can be passed
// to the search function.  These are flags
// an can be bitwise OR'd together to search
// for multiple types of content simultaneously.
const (
	SearchTypeAlbum    SearchType = 1 << iota
	SearchTypeArtist              = 1 << iota
	SearchTypePlaylist            = 1 << iota
	SearchTypeTrack               = 1 << iota
)

func (st SearchType) encode() string {
	types := []string{}
	if st&SearchTypeAlbum != 0 {
		types = append(types, "album")
	}
	if st&SearchTypeArtist != 0 {
		types = append(types, "artist")
	}
	if st&SearchTypePlaylist != 0 {
		types = append(types, "playlist")
	}
	if st&SearchTypeTrack != 0 {
		types = append(types, "track")
	}
	return strings.Join(types, ",")
}

// TODO: maybe instead of exposing the prev/next URLs,
// we can just have functions for retrieving the prev/next page
type resultPage struct {
	// A link to the Web API Endpoint returning the full
	// result of this request.
	FullResult string
	// The maximum number of items in the response, as set
	// in the query (or default value if unset).
	Limit int
	// The offset of the items returned, as set in the query
	// (or default value if unset).
	Offset int
	// The total number of items available to return.
	Total int
	// The URL to the next page of items (if available).
	Next string
	// The URL to the previous page of items (if available).
	Previous string
}

// ArtistResult contains artists returned by the Web API.
type ArtistResult struct {
	resultPage
	Artists []FullArtist
}

// AlbumResult contains albums returned by the Web API.
type AlbumResult struct {
	resultPage
	Albums []SimpleAlbum
}

// PlaylistResult contains playlists returned by the Web API.
type PlaylistResult struct {
	resultPage
	Playlists []SimplePlaylist
}

// TrackResult contains tracks returned by the Web API.
type TrackResult struct {
	resultPage
	Tracks []SimpleTrack
}

type searchResult struct {
	Artists   *page `json:"artists"`
	Albums    *page `json:"albums"`
	Tracks    *page `json:"tracks"`
	Playlists *page `json:"playlists"`
}

// SearchOptions contains optional parameters for the search functions.
// Only the non-nil fields are used in the query.
type SearchOptions struct {
	// The maximum number of objects to return.  If not specified,
	// Spotify will return 20 results by default.  Minimum: 1,
	// Maximum 50.
	Limit *int
	// The index of the first object to return.  If not specified,
	// Spotify returns the first object.  Use with Limit to get the
	// next set of results.
	Offset *int
	// An ISO 3166-1 alpha-2 country code, or the constant string
	// MarketFromToken.  If specified, only artists, albums, and tracks
	// with content playable in the specified market will be returned.
	// (Playlist results are not affected by the market parameter)
	Market *string
}

// SearchResult contains the results of a call to Search.
// Fields that weren't searched for will be nil pointers.
type SearchResult struct {
	Artists   *ArtistResult
	Albums    *AlbumResult
	Playlists *PlaylistResult
	Tracks    *TrackResult
}

// Search is a wrapper around DefaultClient.Search.
func Search(query string, t SearchType) (*SearchResult, error) {
	return DefaultClient.Search(query, t)
}

// SearchFiltered is a wrapper around DefaultClient.SearchFiltered
func SearchFiltered(query string, t SearchType, opt *SearchOptions) (*SearchResult, error) {
	return DefaultClient.SearchFiltered(query, t, opt)
}

// Search gets Spotify catalog information about artists,
// albums, tracks, or playlists that match a keyword string.
// t is a mask containing one or more search types.  For
// example, Search(query, Artist | Album) will search for
// artists or albums matching the specified keywords.
//
// Matching of search keywords is NOT case sensitive.  Keywords
// are matched in any order unless surrounded by double quotes.
// Searching for playlists will return results where the query
// keyword(s) match any part of the playlist's name or description.
// Only popular public playlists are returned.
//
// Operators
//
// The operator NOT can be used to exclude results.  For example,
// query = "roadhouse NOT blues" returns items that match
// roadhouse but exludes those that also contain the keyword
// "blues".  Similarly, the OR operator can be used to broaden
// the search.  query = "roadhouse OR blues" returns all results
// that include either of the terms.  Only one OR operator can
// be used in a query.
//
// Operators should be specified in uppercase.
//
// Wildcards
//
// The asterisk (*) character can, with some limitations, be used
// as a wildcard (maximum of 2 per query).  It will match a
// variable number of non-white-space characters.  It cannot be
// used in a quoted phrase, in a field filter, or as the first
// character of a keyword string.
//
// Field filters
//
// By default, results are returned when a match is found in
// any field of the target object type.  Searches can be made
// more specific by specifying an album, artist, or track
// field filter.  For example, "album:gold artist:abba type:album"
// will only return results with the text "gold" in the album
// name and the text "abba" in the artist's name.
//
// The field filter "year" can be used with album, artist, and
// track searches to limit the results to a particular year.
// For example "bob year:2014" or "bob year:1980-2020".
//
// The field filter "tag:new" can be used in album searches
// to retrieve only albums released in the last two weeks.
// The field filter "tag:hipster" can be used in album
// searches to retrieve only albums with the lowest 10%
// popularity.
//
// Other possible field filters, depending on object types
// being searched, include "genre", "upc", and "isrc".
// For example "damian genre:reggae-pop".
func (c *Client) Search(query string, t SearchType) (*SearchResult, error) {
	return c.SearchFiltered(query, t, nil)
}

// SearchFiltered works just like Search, but it accepts additional
// parameters for filtering the output.
func (c *Client) SearchFiltered(query string, t SearchType, opt *SearchOptions) (*SearchResult, error) {
	query = url.QueryEscape(query)
	v := url.Values{}
	v.Set("q", query)
	v.Set("type", t.encode())
	if opt != nil {
		if opt.Limit != nil {
			v.Set("limit", strconv.Itoa(*opt.Limit))
		}
		if opt.Market != nil {
			v.Set("market", *opt.Market)
		}
		if opt.Offset != nil {
			v.Set("offset", strconv.Itoa(*opt.Offset))
		}
	}
	uri := baseAddress + "search?" + v.Encode()
	resp, err := c.http.Get(uri)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, decodeError(resp.Body)
	}

	var result searchResult
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return nil, err
	}
	sr := &SearchResult{
		Artists:   toArtists(result.Artists),
		Playlists: toPlaylists(result.Playlists),
		Albums:    toAlbums(result.Albums),
		Tracks:    toTracks(result.Tracks),
	}
	return sr, err
}

func toArtists(p *page) *ArtistResult {
	if p == nil {
		return nil
	}
	var a ArtistResult
	a.FullResult = p.Endpoint
	a.Limit = p.Limit
	a.Offset = p.Offset
	a.Total = p.Total
	a.Previous = p.Previous
	a.Next = p.Next

	err := json.Unmarshal([]byte(p.Items), &a.Artists)
	if err != nil {
		return nil
	}
	return &a
}

func toAlbums(p *page) *AlbumResult {
	if p == nil {
		return nil
	}
	var a AlbumResult
	a.FullResult = p.Endpoint
	a.Limit = p.Limit
	a.Offset = p.Offset
	a.Total = p.Total
	a.Previous = p.Previous
	a.Next = p.Next

	err := json.Unmarshal([]byte(p.Items), &a.Albums)
	if err != nil {
		return nil
	}
	return &a
}

func toPlaylists(p *page) *PlaylistResult {
	if p == nil {
		return nil
	}
	var a PlaylistResult
	a.FullResult = p.Endpoint
	a.Limit = p.Limit
	a.Offset = p.Offset
	a.Total = p.Total
	a.Previous = p.Previous
	a.Next = p.Next

	err := json.Unmarshal([]byte(p.Items), &a.Playlists)
	if err != nil {
		return nil
	}
	return &a
}

func toTracks(p *page) *TrackResult {
	if p == nil {
		return nil
	}
	var a TrackResult
	a.FullResult = p.Endpoint
	a.Limit = p.Limit
	a.Offset = p.Offset
	a.Total = p.Total
	a.Previous = p.Previous
	a.Next = p.Next

	err := json.Unmarshal([]byte(p.Items), &a.Tracks)
	if err != nil {
		return nil
	}
	return &a
}
