// Package spotify provides utilities for interfacing
// with Spotify's Web API.
package spotify

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"golang.org/x/oauth2"
)

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
)

// Client is a client for working with the Spotify Web API.
// It is best to create this using spotify.New()
type Client struct {
	http    *http.Client
	baseURL string

	autoRetry      bool
	acceptLanguage string
}

type ClientOption func(client *Client)

// WithRetry configures the Spotify API client to automatically retry requests that fail due to rate limiting.
func WithRetry(shouldRetry bool) ClientOption {
	return func(client *Client) {
		client.autoRetry = shouldRetry
	}
}

// WithBaseURL provides an alternative base url to use for requests to the Spotify API. This can be used to connect to a
// staging or other alternative environment.
func WithBaseURL(url string) ClientOption {
	return func(client *Client) {
		client.baseURL = url
	}
}

// WithAcceptLanguage configures the client to provide the accept language header on all requests.
func WithAcceptLanguage(lang string) ClientOption {
	return func(client *Client) {
		client.acceptLanguage = lang
	}
}

// New returns a client for working with the Spotify Web API.
// The provided httpClient must provide Authentication with the requests.
// The auth package may be used to generate a suitable client.
func New(httpClient *http.Client, opts ...ClientOption) *Client {
	c := &Client{
		http:    httpClient,
		baseURL: "https://api.spotify.com/v1/",
	}

	for _, opt := range opts {
		opt(c)
	}

	return c
}

// URI identifies an artist, album, track, or category.  For example,
// spotify:track:6rqhFgbbKwnb9MLmUQDhG6
type URI string

// ID is a base-62 identifier for an artist, track, album, etc.
// It can be found at the end of a spotify.URI.
type ID string

func (id *ID) String() string {
	return string(*id)
}

// Numeric is a convenience type for handling numbers sent as either integers or floats.
type Numeric int

// UnmarshalJSON unmarshals a JSON number (float or int) into the Numeric type.
func (n *Numeric) UnmarshalJSON(data []byte) error {
	var f float64
	if err := json.Unmarshal(data, &f); err != nil {
		return err
	}
	*n = Numeric(int(f))
	return nil
}

// Followers contains information about the number of people following a
// particular artist or playlist.
type Followers struct {
	// The total number of followers.
	Count Numeric `json:"total"`
	// A link to the Web API endpoint providing full details of the followers,
	// or the empty string if this data is not available.
	Endpoint string `json:"href"`
}

// Image identifies an image associated with an item.
type Image struct {
	// The image height, in pixels.
	Height Numeric `json:"height"`
	// The image width, in pixels.
	Width Numeric `json:"width"`
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
	// RetryAfter contains the time before which client should not retry a
	// rate-limited request, calculated from the Retry-After header, when present.
	RetryAfter time.Time `json:"-"`
}

func (e Error) Error() string {
	return fmt.Sprintf("spotify: %s [%d]", e.Message, e.Status)
}

// HTTPStatus returns the HTTP status code returned by the server when the error
// occurred.
func (e Error) HTTPStatus() int {
	return e.Status
}

// decodeError decodes an Error from an io.Reader.
func decodeError(resp *http.Response) error {
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	if ctHeader := resp.Header.Get("Content-Type"); ctHeader == "" {
		msg := string(responseBody)
		if len(msg) == 0 {
			msg = http.StatusText(resp.StatusCode)
		}

		return Error{
			Message: msg,
			Status:  resp.StatusCode,
		}
	}

	if len(responseBody) == 0 {
		return Error{
			Message: "server response without body",
			Status:  resp.StatusCode,
		}
	}

	buf := bytes.NewBuffer(responseBody)

	var e struct {
		E Error `json:"error"`
	}
	err = json.NewDecoder(buf).Decode(&e)
	if err != nil {
		return Error{
			Message: fmt.Sprintf("failed to decode error response %q", responseBody),
			Status:  resp.StatusCode,
		}
	}

	e.E.Status = resp.StatusCode
	if e.E.Message == "" {
		// Some errors will result in there being a useful status-code but an
		// empty message. An example of this is when we send some of the
		// arguments directly in the HTTP query and the URL ends-up being too
		// long.

		e.E.Message = "server response without error description"
	}
	if retryAfter, _ := strconv.Atoi(resp.Header.Get("Retry-After")); retryAfter != 0 {
		e.E.RetryAfter = time.Now().Add(time.Duration(retryAfter) * time.Second)
	}

	return e.E
}

// shouldRetry determines whether the status code indicates that the
// previous operation should be retried at a later time
func shouldRetry(status int) bool {
	return status == http.StatusAccepted || status == http.StatusTooManyRequests
}

// isFailure determines whether the code indicates failure
func isFailure(code int, validCodes []int) bool {
	for _, item := range validCodes {
		if item == code {
			return false
		}
	}
	return true
}

// `execute` executes a non-GET request. `needsStatus` describes other HTTP
// status codes that will be treated as success. Note that we allow all 200s
// even if there are additional success codes that represent success.
func (c *Client) execute(req *http.Request, result interface{}, needsStatus ...int) error {
	if c.acceptLanguage != "" {
		req.Header.Set("Accept-Language", c.acceptLanguage)
	}
	for {
		resp, err := c.http.Do(req)
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		if c.autoRetry &&
			isFailure(resp.StatusCode, needsStatus) &&
			shouldRetry(resp.StatusCode) {
			select {
			case <-req.Context().Done():
				// If the context is cancelled, return the original error
			case <-time.After(retryDuration(resp)):
				continue
			}
		}
		if resp.StatusCode == http.StatusNoContent {
			return nil
		}
		if (resp.StatusCode >= 300 ||
			resp.StatusCode < 200) &&
			isFailure(resp.StatusCode, needsStatus) {
			return decodeError(resp)
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

func retryDuration(resp *http.Response) time.Duration {
	raw := resp.Header.Get("Retry-After")
	if raw == "" {
		return defaultRetryDuration
	}
	seconds, err := strconv.ParseInt(raw, 10, 32)
	if err != nil {
		return defaultRetryDuration
	}
	return time.Duration(seconds) * time.Second
}

func (c *Client) get(ctx context.Context, url string, result interface{}) error {
	for {
		req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
		if c.acceptLanguage != "" {
			req.Header.Set("Accept-Language", c.acceptLanguage)
		}
		if err != nil {
			return err
		}
		resp, err := c.http.Do(req)
		if err != nil {
			return err
		}

		defer resp.Body.Close()

		if resp.StatusCode == http.StatusTooManyRequests && c.autoRetry {
			select {
			case <-ctx.Done():
				// If the context is cancelled, return the original error
			case <-time.After(retryDuration(resp)):
				continue
			}
		}
		if resp.StatusCode == http.StatusNoContent {
			return nil
		}
		if resp.StatusCode != http.StatusOK {
			return decodeError(resp)
		}

		return json.NewDecoder(resp.Body).Decode(result)
	}
}

// Token gets the client's current token.
func (c *Client) Token() (*oauth2.Token, error) {
	transport, ok := c.http.Transport.(*oauth2.Transport)
	if !ok {
		return nil, errors.New("spotify: client not backed by oauth2 transport")
	}
	t, err := transport.Source.Token()
	if err != nil {
		return nil, err
	}
	return t, nil
}
