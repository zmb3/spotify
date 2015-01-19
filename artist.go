package spotify

import (
	"encoding/json"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

// SimpleArtist contains basic info about an artist.
type SimpleArtist struct {
	// The name of the artist.
	Name string `json:"name"`
	// The Spotify ID for the artist.
	ID ID `json:"id"`
	// The Spotify URI for the artist.
	URI URI `json:"uri"`
	// A link to the Web API enpoint providing
	// full details of the artist.
	Endpoint string `json:"href"`
	// Known external URLs for this artist.
	ExternalURLs ExternalURL `json:"external_urls"`
}

// FullArtist provides extra artist data in addition
// to what is provided by SimpleArtist.
type FullArtist struct {
	SimpleArtist
	// The popularity of the artist.  The value will be
	// between 0 and 100, with 100 being the most popular.
	// The artist's popularity is calculated from the
	// popularity of all of the artist's tracks.
	Popularity int `json:"popularity"`
	// A list of genres the artist is associated with.
	// For example, "Prog Rock" or "Post-Grunge".  If
	// not yet classified, the slice is empty.
	Genres []string `json:"genres"`
	// Information about followers of the artist.
	Followers Followers
	// Images of the artist in various sizes, widest first.
	Images []Image `json:"images"`
}

// FindArtist is a wrapper around DefaultClient.FindArtist.
func FindArtist(id ID) (*FullArtist, error) {
	return DefaultClient.FindArtist(id)
}

// FindArtist gets Spotify catalog information for a single
// artist, given that artist's Spotify ID.
func (c *Client) FindArtist(id ID) (*FullArtist, error) {
	uri := baseAddress + "artists/" + string(id)
	resp, err := c.http.Get(uri)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, decodeError(resp.Body)
	}
	var a FullArtist
	err = json.NewDecoder(resp.Body).Decode(&a)
	if err != nil {
		return nil, err
	}
	return &a, nil
}

// FindArtists is a wrapper around DefaultClient.FindArtists.
func FindArtists(ids ...ID) ([]*FullArtist, error) {
	return DefaultClient.FindArtists(ids...)
}

// FindArtists gets spotify catalog information for several
// artists based on their Spotify IDs.  It supports up to
// 50 artists in a single call.  Artists are returned in the
// order requested.  If an artist is not found, that position
// in the result will be nil.  Duplicate IDs will result in
// duplicate artists in the result.
func (c *Client) FindArtists(ids ...ID) ([]*FullArtist, error) {
	uri := baseAddress + "artists?ids=" + strings.Join(toStringSlice(ids), ",")
	resp, err := c.http.Get(uri)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, decodeError(resp.Body)
	}
	var a struct {
		Artists []*FullArtist
	}
	err = json.NewDecoder(resp.Body).Decode(&a)
	if err != nil {
		return nil, err
	}
	return a.Artists, nil
}

// ArtistsTopTracks is a wrapper around DefaultClient.ArtistTopTracks.
func ArtistsTopTracks(artistID ID, country string) ([]FullTrack, error) {
	return DefaultClient.ArtistsTopTracks(artistID, country)
}

// ArtistsTopTracks gets Spotify catalog information about
// an artist's top tracks in a particular country.  It returns
// a maximum of 10 tracks.  The country is specified as an
// ISO 3166-1 alpha-2 country code.
func (c *Client) ArtistsTopTracks(artistID ID, country string) ([]FullTrack, error) {
	uri := baseAddress + "artists/" + string(artistID) + "/top-tracks?country=" + country
	resp, err := c.http.Get(uri)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, decodeError(resp.Body)
	}
	var t struct {
		Tracks []FullTrack `json:"tracks"`
	}

	err = json.NewDecoder(resp.Body).Decode(&t)
	if err != nil {
		return nil, err
	}
	return t.Tracks, nil
}

// FindRelatedArtists is a wrapper around DefaultClient.FindRelatedArtists.
func FindRelatedArtists(id ID) ([]FullArtist, error) {
	return DefaultClient.FindRelatedArtists(id)
}

// FindRelatedArtists gets Spotify catalog information about
// artists similar to a given artist.  Similarity is based on
// analysis of the Spotify community's listening history.
// This function returns up to 20 artists that are considered
// related to the specified artist.
func (c *Client) FindRelatedArtists(id ID) ([]FullArtist, error) {
	uri := baseAddress + "artists/" + string(id) + "/related-artists"
	resp, err := c.http.Get(uri)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, decodeError(resp.Body)
	}
	var a struct {
		Artists []FullArtist `json:"artists"`
	}
	err = json.NewDecoder(resp.Body).Decode(&a)
	if err != nil {
		return nil, err
	}
	return a.Artists, nil
}

// ArtistAlbums is a wrapper around DefaultClient.ArtistAlbums.
func ArtistAlbums(artistID ID) (*AlbumResult, error) {
	return DefaultClient.ArtistAlbums(artistID)
}

// ArtistAlbums gets Spotify catalog information
// about an artist's albums.  It is equivalent to
// ArtistAlbumsFiltered(artistID, nil)
func (c *Client) ArtistAlbums(artistID ID) (*AlbumResult, error) {
	return c.ArtistAlbumsFiltered(artistID, nil)
}

// AlbumOptions contains optional parameters for sorting
// and filtering the results of FindArtistAlbumsFiltered.
// Only the non-nil fields are used in the query.
type AlbumOptions struct {
	// Type contains one or more types of albums
	// that will be used to filter the result.  If
	// not specified, all album types will be returned.
	Type *AlbumType
	// An ISO 3166-1 alpha-2 country code that  limits
	// the result to one particular market.  If not given,
	// results will be returned for all markets and you are
	// likely to get duplicate results per album, one for
	// each market in which the album is available!
	Country *string
	// The maximum number of objects to return.
	// Minimum: 1.  Maximum: 50.
	Limit *int
	// The index of the first album to return. Use
	// with Limit to get the next set of albums.
	Offset *int
}

// SetType sets the optional AlbumType field.
func (a *AlbumOptions) SetType(t AlbumType) {
	a.Type = new(AlbumType)
	*a.Type = t
}

// SetCountry sets the optional Country field.
func (a *AlbumOptions) SetCountry(c string) {
	a.Country = new(string)
	*a.Country = c
}

// SetLimit sets the optional Limit field.
func (a *AlbumOptions) SetLimit(l int) {
	a.Limit = new(int)
	*a.Limit = l
}

// SetOffset sets the optional Offset field.
func (a *AlbumOptions) SetOffset(o int) {
	a.Offset = new(int)
	*a.Offset = o
}

// ArtistAlbumsFiltered is a wrapper around DefaultClient.ArtistAlbumsFiltered
func ArtistAlbumsFiltered(artistID ID, options *AlbumOptions) (*AlbumResult, error) {
	return DefaultClient.ArtistAlbumsFiltered(artistID, options)
}

// ArtistAlbumsFiltered is just like ArtistAlbums, but
// it accepts optional parameters to filter and sort the result.
func (c *Client) ArtistAlbumsFiltered(artistID ID, options *AlbumOptions) (*AlbumResult, error) {
	uri := baseAddress + "artists/" + string(artistID) + "/albums"
	// add optional query string if options were specified
	if options != nil {
		values := url.Values{}
		if options.Type != nil {
			values.Set("album_type", options.Type.encode())
		}
		if options.Country != nil {
			values.Set("market", *options.Country)
		}
		if options.Limit != nil {
			values.Set("limit", strconv.Itoa(*options.Limit))
		}
		if options.Offset != nil {
			values.Set("offset", strconv.Itoa(*options.Offset))
		}
		if query := values.Encode(); query != "" {
			uri += "?" + query
		}
	}
	resp, err := c.http.Get(uri)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, decodeError(resp.Body)
	}
	var p page
	err = json.NewDecoder(resp.Body).Decode(&p)
	if err != nil {
		return nil, err
	}
	var result AlbumResult
	err = json.Unmarshal([]byte(p.Items), &result.Albums)
	if err != nil {
		return nil, err
	}
	result.FullResult = p.Endpoint
	result.Limit = p.Limit
	result.Offset = p.Offset
	result.Total = p.Total
	result.Previous = p.Previous
	result.Next = p.Next
	return &result, nil
}
