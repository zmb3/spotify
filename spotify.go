// Package spotify provides utilties for interfacing
// with Spotify's Web API.
package spotify

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

// Version is the version of this library.
const Version = "1.0.0"

const (
	// DateLayout can be used with time.Parse to create time.Time values
	// from Spotify date strings.  For example, PrivateUser.Birthdate
	// uses this format.
	DateLayout = "2006-01-02"
	// TimestampLayout can be used with time.Parse to create time.Time
	// values from SpotifyTimestamp strings.  It is an ISO 8601 UTC timestamp
	// with a zero offset.  For example, PlaylistTrack's AddedAt field uses
	// this format.
	TimestampLayout = "2006-01-02T15:04:05Z"

	// defaultRetryDurationS helps us fix an apparent server bug whereby we will
	// be told to retry but not be given a wait-interval.
	defaultRetryDuration = time.Second * 5

	// rateLimitExceededStatusCode is the code that the server returns when our
	// request frequency is too high.
	rateLimitExceededStatusCode = 429
)

var (
	baseAddress = "https://api.spotify.com/v1/"

	// DefaultClient is the default client that is used by the wrapper functions
	// that don't require authorization.  If you need to authenticate, create
	// your own client with `Authenticator.NewClient`.
	DefaultClient = &Client{
		http: new(http.Client),
	}
)

// URI identifies an artist, album, track, or category.  For example,
// spotify:track:6rqhFgbbKwnb9MLmUQDhG6
type URI string

// ID is a base-62 identifier for an artist, track, album, etc.
// It can be found at the end of a spotify.URI.
type ID string

func init() {
	// disable HTTP/2 for DefaultClient, see: https://github.com/zmb3/spotify/issues/20
	tr := &http.Transport{
		TLSNextProto: map[string]func(authority string, c *tls.Conn) http.RoundTripper{},
	}
	DefaultClient.http.Transport = tr
}

func (id *ID) String() string {
	return string(*id)
}

// Followers contains information about the number of people following a
// particular artist or playlist.
type Followers struct {
	// The total number of followers.
	Count uint `json:"total"`
	// A link to the Web API endpoint providing full details of the followers,
	// or the empty string if this data is not available.
	Endpoint string `json:"href"`
}

// Image identifies an image associated with an item.
type Image struct {
	// The image height, in pixels.
	Height int `json:"height"`
	// The image width, in pixels.
	Width int `json:"width"`
	// The source URL of the image.
	URL string `json:"url"`
}

// Download downloads the image and writes its data to the specified io.Writer.
func (i Image) Download(dst io.Writer) error {
	resp, err := http.Get(i.URL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	// TODO: get Content-Type from header?
	if resp.StatusCode != http.StatusOK {
		return errors.New("Couldn't download image - HTTP" + strconv.Itoa(resp.StatusCode))
	}
	_, err = io.Copy(dst, resp.Body)
	return err
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

// decodeError decodes an Error from an io.Reader.
func (c *Client) decodeError(resp *http.Response) error {
	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if len(responseBody) == 0 {
		return fmt.Errorf("spotify: HTTP %d: %s (body empty)", resp.StatusCode, http.StatusText(resp.StatusCode))
	}

	buf := bytes.NewBuffer(responseBody)

	var e struct {
		E Error `json:"error"`
	}
	err = json.NewDecoder(buf).Decode(&e)
	if err != nil {
		return fmt.Errorf("spotify: couldn't decode error: (%d) [%s]", len(responseBody), responseBody)
	}

	if resp.StatusCode == rateLimitExceededStatusCode {
		retrySecondsRaw := resp.Header.Get("Retry-After")
		if retrySecondsRaw != "" {
			retrySeconds, err := strconv.ParseInt(retrySecondsRaw, 10, 32)
			if err != nil {
				return fmt.Errorf("spotify: could not parse retry seconds: %s", retrySecondsRaw)
			}

			if retrySeconds == 0 {
				c.retryDuration = defaultRetryDuration
			} else {
				c.retryDuration = time.Second * time.Duration(retrySeconds)
			}
		}
	} else if e.E.Error() == "" {
		// Some errors will result in there being a useful status-code but an
		// empty message, which will confuse the user (who only has access to
		// the message and not the code). An example of this is when we send
		// some of the arguments directly in the HTTP query and the URL ends-up
		// being too long.

		e.E.Message = fmt.Sprintf("spotify: unexpected HTTP %d: %s (empty error)", resp.StatusCode, http.StatusText(resp.StatusCode))
	}

	return e.E
}

// Client is a client for working with the Spotify Web API.
// To create an authenticated client, use the
// `Authenticator.NewClient` method.  If you don't need to
// authenticate, you can use `DefaultClient`.
type Client struct {
	http *http.Client

	AutoRetry     bool
	retryDuration time.Duration
}

// executeOpt executes a non-GET request. `needsStatus` describes another code
// that can represent success. Note that in all current usages of this function,
// we need to still allow a 200 even if we'd also like to check for a second
// success code.
func (c *Client) executeOpt(req *http.Request, needsStatus int, result interface{}) error {
	for {
		resp, err := c.http.Do(req)
		if err != nil {
			return err
		}

		defer resp.Body.Close()

		if resp.StatusCode == rateLimitExceededStatusCode && c.AutoRetry {
			time.Sleep(c.retryDuration)
			continue
		} else if resp.StatusCode != http.StatusOK && (needsStatus == 0 || resp.StatusCode != needsStatus) {
			errorMessage := c.decodeError(resp)
			return errorMessage
		}

		if result != nil {
			if err := json.NewDecoder(resp.Body).Decode(result); err != nil {
				return err
			}
		}

		break
	}

	return nil
}

func (c *Client) execute(req *http.Request) error {
	return c.executeOpt(req, 0, nil)
}

func (c *Client) get(url string, result interface{}) error {
	for {
		resp, err := c.http.Get(url)
		if err != nil {
			return err
		}

		defer resp.Body.Close()

		if resp.StatusCode == rateLimitExceededStatusCode && c.AutoRetry {
			time.Sleep(c.retryDuration)
			continue
		} else if resp.StatusCode != http.StatusOK {
			errorMessage := c.decodeError(resp)
			return errorMessage
		}

		err = json.NewDecoder(resp.Body).Decode(result)
		if err != nil {
			return err
		}

		break
	}

	return nil
}

// Options contains optional parameters that can be provided
// to various API calls.  Only the non-nil fields are used
// in queries.
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

// NewReleasesOpt is like NewReleases, but it accepts optional parameters
// for filtering the results.
func (c *Client) NewReleasesOpt(opt *Options) (albums *SimpleAlbumPage, err error) {
	spotifyURL := baseAddress + "browse/new-releases"
	if opt != nil {
		v := url.Values{}
		if opt.Country != nil {
			v.Set("country", *opt.Country)
		}
		if opt.Limit != nil {
			v.Set("limit", strconv.Itoa(*opt.Limit))
		}
		if opt.Offset != nil {
			v.Set("offset", strconv.Itoa(*opt.Offset))
		}
		if params := v.Encode(); params != "" {
			spotifyURL += "?" + params
		}
	}

	var objmap map[string]*json.RawMessage
	err = c.get(spotifyURL, &objmap)
	if err != nil {
		return nil, err
	}

	var result SimpleAlbumPage
	err = json.Unmarshal(*objmap["albums"], &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

// NewReleases gets a list of new album releases featured in Spotify.
// This call requires bearer authorization.
func (c *Client) NewReleases() (albums *SimpleAlbumPage, err error) {
	return c.NewReleasesOpt(nil)
}
