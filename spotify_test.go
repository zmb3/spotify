package spotify

import (
	"errors"
	"net/http"
	"os"
	"strings"
	"testing"
)

type stringRoundTripper struct {
	strings.Reader
	statusCode  int
	lastRequest *http.Request
}

func newStringRoundTripper(code int, s string) *stringRoundTripper {
	return &stringRoundTripper{*strings.NewReader(s), code, nil}
}

func (s stringRoundTripper) Close() error {
	return nil
}

type fileRoundTripper struct {
	*os.File
	statusCode  int
	lastRequest *http.Request
}

func newFileRoundTripper(code int, filename string) *fileRoundTripper {
	file, err := os.Open(filename)
	if err != nil {
		panic("Couldn't open file " + filename)
	}
	return &fileRoundTripper{file, code, nil}
}

func (s *stringRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	s.lastRequest = req
	if req.Header == nil {
		if req.Body != nil {
			req.Body.Close()
		}
		return nil, errors.New("stringRoundTripper: nil request header")
	}
	return &http.Response{
		StatusCode: s.statusCode,
		Body:       s,
	}, nil
}

func (f *fileRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	f.lastRequest = req
	if req.Header == nil {
		if req.Body != nil {
			req.Body.Close()
		}
		return nil, errors.New("fileRoundTripper: nil request header")
	}
	return &http.Response{
		StatusCode: f.statusCode,
		Body:       f,
	}, nil
}

// Returns a client whose requests will always return
// the specified status code and body.
func testClientString(code int, body string) *Client {
	return &Client{
		http: &http.Client{
			Transport: newStringRoundTripper(code, body),
		},
	}
}

// Returns a client whose requests will always return
// a response with the specified status code and a body
// that is read from the specified file.
func testClientFile(code int, filename string) *Client {
	return &Client{
		http: &http.Client{
			Transport: newFileRoundTripper(code, filename),
		},
	}
}

func getLastRequest(c *Client) *http.Request {
	if frt, ok := c.http.Transport.(*fileRoundTripper); ok {
		return frt.lastRequest
	}
	if srt, ok := c.http.Transport.(*stringRoundTripper); ok {
		return srt.lastRequest
	}
	return nil
}

// addDummyAuth puts fake authorization data in the specified
// client, which allows the basic authentication checks to pass
// for the purpose of testing
func addDummyAuth(c *Client) {
	// c.AccessToken = "sample token"
	// c.TokenType = BearerToken
}

func TestNewReleases(t *testing.T) {
	c := testClientFile(http.StatusOK, "test_data/new_releases.txt")
	addDummyAuth(c)
	r, err := c.NewReleases()
	if err != nil {
		t.Error(err)
		return
	}
	if r.Albums[0].ID != "60mvULtYiNSRmpVvoa3RE4" {
		t.Error("Invalid data: ", r.Albums[0].ID)
		return
	}
	if r.Albums[0].Name != "We Are One (Ole Ola) [The Official 2014 FIFA World Cup Song]" {
		t.Error("Invalid data", r.Albums[0].Name)
		return
	}
}
