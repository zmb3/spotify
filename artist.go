package spotify

import (
	"encoding/json"
	"errors"
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

// FindArtist gets Spotify catalog information for a single
// artist, given that artist's Spotify ID.
func (c *Client) FindArtist(id ID) (*FullArtist, error) {
	uri := baseAddress + "artists/" + string(id)
	resp, err := c.http.Get(uri)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		var e struct {
			E Error `json:"error"`
		}
		err = json.NewDecoder(resp.Body).Decode(&e)
		if err != nil {
			return nil, errors.New("spotify: HTTP response error")
		}
		return nil, e.E
	}
	var a FullArtist
	err = json.NewDecoder(resp.Body).Decode(&a)
	if err != nil {
		return nil, err
	}
	return &a, nil
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
	if resp.StatusCode != 200 {
		var e struct {
			E Error `json:"error"`
		}
		err = json.NewDecoder(resp.Body).Decode(&e)
		if err != nil {
			return nil, errors.New("spotify: HTTP response error")
		}
		return nil, e.E
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
