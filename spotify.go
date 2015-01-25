// Package spotify provides utilties for interfacing
// with Spotify's Web API.
package spotify

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

var (
	baseAddress = "https://api.spotify.com/v1/"

	// DefaultClient is the default client and is used by ...
	DefaultClient = &Client{}

	// ErrNotAuthenticated is returned when an unauthenticated user
	// makes an API call that requries authentication.
	ErrNotAuthenticated = errors.New("spotify: this call requires authentication")
)

// URI identifies an artist, album, or track.  For example,
// spotify:track:6rqhFgbbKwnb9MLmUQDhG6
type URI string

// ID is a base-62 identifier for an artist, track, album, etc.
// It can be found at the end of a spotify.URI.
type ID string

func (id *ID) String() string {
	return string(*id)
}

// Timestamp is an ISO 8601 formatted timestamp
// representing Coordinated Universal Time (UTC)
// with zero offset: YYYY-MM-DDTHH:MM:SSZ.
type Timestamp string

// Followers contains information about the number of people following a
// particular artist or playlist.
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

// Error represents an error returned by the Spotify Web API.
type Error struct {
	// A short description of the error.
	Message string `json:"message"`
	// The HTTP status code.
	Status int `json:"status"`
}

func (e Error) Error() string {
	return e.Message
}

// decodeError decodes an error from an io.Reader.
func decodeError(r io.Reader) error {
	var e struct {
		E Error `json:"error"`
	}
	err := json.NewDecoder(r).Decode(&e)
	if err != nil {
		return errors.New("spotify: couldn't decode error")
	}
	return e.E
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
	http        http.Client
	AccessToken string
	TokenType   TokenType
}

// Options contains optional parameters that can be provided
// to various API calls.  Only the non-nil fields are used
// in queries.
//
//
type Options struct {
	// Country is an ISO 3166-1 alpha-2 country code.  Provide
	// this parameter if you want the list of returned items to
	// be relevant to a particular country.  If omitted, the
	// results will be relevant to all countries.
	Country *string
	// Limit is the maximum number of items to return.
	Limit *int
	// Offset is the index of the first item to return.  Use it
	// with Limit to get the next set of items.
	Offset *int
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
