package spotify

import (
	"context"
	"fmt"
	"strings"
)

// SimpleArtist contains basic info about an artist.
type SimpleArtist struct {
	Name string `json:"name"`
	ID   ID     `json:"id"`
	// The Spotify URI for the artist.
	URI URI `json:"uri"`
	// A link to the Web API enpoint providing full details of the artist.
	Endpoint     string            `json:"href"`
	ExternalURLs map[string]string `json:"external_urls"`
}

// FullArtist provides extra artist data in addition to what is provided by SimpleArtist.
type FullArtist struct {
	SimpleArtist
	// The popularity of the artist, expressed as an integer between 0 and 100.
	// The artist's popularity is calculated from the popularity of the artist's tracks.
	Popularity int `json:"popularity"`
	// A list of genres the artist is associated with.  For example, "Prog Rock"
	// or "Post-Grunge".  If not yet classified, the slice is empty.
	Genres    []string  `json:"genres"`
	Followers Followers `json:"followers"`
	// Images of the artist in various sizes, widest first.
	Images []Image `json:"images"`
}

// GetArtist gets Spotify catalog information for a single artist, given its Spotify ID.
func (c *Client) GetArtist(ctx context.Context, id ID) (*FullArtist, error) {
	spotifyURL := fmt.Sprintf("%sartists/%s", c.baseURL, id)

	var a FullArtist
	err := c.get(ctx, spotifyURL, &a)
	if err != nil {
		return nil, err
	}

	return &a, nil
}

// GetArtists gets spotify catalog information for several artists based on their
// Spotify IDs.  It supports up to 50 artists in a single call.  Artists are
// returned in the order requested.  If an artist is not found, that position
// in the result will be nil.  Duplicate IDs will result in duplicate artists
// in the result.
func (c *Client) GetArtists(ctx context.Context, ids ...ID) ([]*FullArtist, error) {
	spotifyURL := fmt.Sprintf("%sartists?ids=%s", c.baseURL, strings.Join(toStringSlice(ids), ","))

	var a struct {
		Artists []*FullArtist
	}

	err := c.get(ctx, spotifyURL, &a)
	if err != nil {
		return nil, err
	}

	return a.Artists, nil
}

// GetArtistsTopTracks gets Spotify catalog information about an artist's top
// tracks in a particular country.  It returns a maximum of 10 tracks.  The
// country is specified as an ISO 3166-1 alpha-2 country code.
func (c *Client) GetArtistsTopTracks(ctx context.Context, artistID ID, country string) ([]FullTrack, error) {
	spotifyURL := fmt.Sprintf("%sartists/%s/top-tracks?country=%s", c.baseURL, artistID, country)

	var t struct {
		Tracks []FullTrack `json:"tracks"`
	}

	err := c.get(ctx, spotifyURL, &t)
	if err != nil {
		return nil, err
	}

	return t.Tracks, nil
}

// GetRelatedArtists gets Spotify catalog information about artists similar to a
// given artist.  Similarity is based on analysis of the Spotify community's
// listening history.  This function returns up to 20 artists that are considered
// related to the specified artist.
func (c *Client) GetRelatedArtists(ctx context.Context, id ID) ([]FullArtist, error) {
	spotifyURL := fmt.Sprintf("%sartists/%s/related-artists", c.baseURL, id)

	var a struct {
		Artists []FullArtist `json:"artists"`
	}

	err := c.get(ctx, spotifyURL, &a)
	if err != nil {
		return nil, err
	}

	return a.Artists, nil
}

// GetArtistAlbums gets Spotify catalog information about an artist's albums.
// It is equivalent to GetArtistAlbumsOpt(artistID, nil).
//
// The AlbumType argument can be used to find a particular types of album.
// If the Market is not specified, Spotify will likely return a lot
// of duplicates (one for each market in which the album is available
//
// Supported options: Market
func (c *Client) GetArtistAlbums(ctx context.Context, artistID ID, ts []AlbumType, opts ...RequestOption) (*SimpleAlbumPage, error) {
	spotifyURL := fmt.Sprintf("%sartists/%s/albums", c.baseURL, artistID)
	// add optional query string if options were specified
	values := processOptions(opts...).urlParams

	if ts != nil {
		types := make([]string, len(ts))
		for i := range ts {
			types[i] = ts[i].encode()
		}
		values.Set("include_groups", strings.Join(types, ","))
	}

	if query := values.Encode(); query != "" {
		spotifyURL += "?" + query
	}

	var p SimpleAlbumPage

	err := c.get(ctx, spotifyURL, &p)
	if err != nil {
		return nil, err
	}

	return &p, nil
}
