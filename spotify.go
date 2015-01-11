// Package spotify provides utilties for interfacing
// with Spotify's Web API.
package spotify

import (
	"encoding/json"
	"net/http"
	"strings"
)

var (
	baseAddress = "https://api.spotify.com/v1/"
)

// SearchType represents the type of a query used
// in the Search function.
type SearchType int

// Search type values that can be passed
// to the search function.  These are flags
// an can be bitwise OR'd together to search
// for multiple types of content simultaneously.
const (
	Album    SearchType = 1 << iota
	Artist              = 1 << iota
	Playlist            = 1 << iota
	Track               = 1 << iota
)

// ISO 3166-1 alpha 2 country codes.
const (
	CountryAustralia = "AU"
	CountryBelgium   = "BE"
	CountryBrazil    = "BR"
	CountryUSA       = "US"
)

// URI identifies an artist,
// album, or track.  For example,
// spotify:track:6rqhFgbbKwnb9MLmUQDhG6
type URI string

// ID is a base-62 identifier for
// an artist, track, album, etc.
// It can be found at the end of
// a spotify.URI.
type ID string

func (id *ID) String() string {
	return string(*id)
}

// Timestamp is an ISO 8601 formatted timestamp
// representing Coordinated Universal Time (UTC)
// with zero offset: YYYY-MM-DDTHH:MM:SSZ.
type Timestamp string

// Followers contains information about the
// number of people following a particular
// artist or playlist.
type Followers struct {
	// The total number of followers.
	Count uint `json:"total"`
	// A link to the Web API endpoint providing
	// full details of the followers, or the empty
	// string if this data is not available.
	Endpoint string `json:"href"`
}

// Image identifies an image associated with an item.
type Image struct {
	// The image height, in pixels.  TODO if unknown?
	Height int `json:"height"`
	// The image width, in pixels.  TODO if unknown?
	Width int `json:"width"`
	// The source URL of the image.
	URL string `json:"url"`
}

// Error represents an error returned by the
// Spotify Web API.
type Error struct {
	// A short description of the error.
	Message string `json:"message"`
	// The HTTP status code.
	Status int `json:"status"`
}

func (e Error) Error() string {
	return e.Message
}

// ExternalID contains information that identifies
// an item.
type ExternalID struct {
	// The identifier type, for example:
	//   "isrc" - International Standard Recording Code
	//   "ean"  - International Article Number
	//   "upc"  - Universal Product Code
	Key string `json:"{key}"`
	// An external identifier for the object.
	Value string `json:"{value}"`
}

// ExternalURL indicates an external, public URL
// for an item.
type ExternalURL struct {
	// The type of the URL, for example:
	//    "spotify" - The Spotify URL for the object.
	Key string `json:"{key}"`
	// An external, public URL to the object.
	Value string `json:"{value}"`
}

// Client is a client for working with the Spotify Web API.
type Client struct {
	http http.Client
}

// NewReleases gets a list of newly released albums that
// are featured in Spotify.
func (c *Client) NewReleases(country string) { // TODO limit/offset
	// get("browse/new-releases")
}

// FeaturedPlaylists gets a list of featured playlists on Spotify.
// This call requires authentication.  The country, locale,
// and timestamp parameters are all optional - pass the
// empty string if you don't care about them.  The country
// parameter allows you to specify an ISO 3166-1 alpha-2
// country code to get featured playlists relevant in a
// particular country.  The locale parameter specifies the
// desired language of the response - it is a lowercase
// ISO 639 language code and an uppercase ISO 3166-1 alpha-2
// country code, joined by an underscore (ie en_US).
// If locale is not supplied, Spotify's default locale is
// used (American English).  The timestamp parameter allows you
// to specify a user's local time (in ISO 8601 format) to
// get results tailored to that specific time.
func (c *Client) FeaturedPlaylists(country, locale, timestamp string) {
	// TODO limit/offset
	// add auth headers
	// header['Content-Type'] = 'application/json';

}

func (st SearchType) encode() string {
	types := []string{}
	if st&Album != 0 {
		types = append(types, "album")
	}
	if st&Artist != 0 {
		types = append(types, "artist")
	}
	if st&Playlist != 0 {
		types = append(types, "playlist")
	}
	if st&Track != 0 {
		types = append(types, "track")
	}
	return strings.Join(types, ",")
}

// page is a container for a set of objects.
// We don't expose this to the user because
// the Items field is just raw JSON.  Instead,
// the user gets AlbumResult, ArtistResult,
// TrackResult, and PlaylistResult.  These
// type all contain the same data as page
// but the Items field is a strongly typed
// slice.
type page struct {
	Endpoint string          `json:"href"`
	Items    json.RawMessage `json:"items"`
	Limit    int             `json:"limit"`
	Next     string          `json:"next"`
	Offset   int             `json:"offset"`
	Previous string          `json:"previous"`
	Total    int             `json:"total"`
}
