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
	statusCode int
}

func newStringRoundTripper(code int, s string) *stringRoundTripper {
	return &stringRoundTripper{*strings.NewReader(s), code}
}

func (s stringRoundTripper) Close() error {
	return nil
}

type fileRoundTripper struct {
	*os.File
	statusCode int
}

func newFileRoundTripper(code int, filename string) *fileRoundTripper {
	file, err := os.Open(filename)
	if err != nil {
		panic("Couldn't open file " + filename)
	}
	return &fileRoundTripper{file, code}
}

func (s *stringRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
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
		http: http.Client{
			Transport: newStringRoundTripper(code, body),
		},
	}
}

// Returns a client whose requests will always return
// a response with the specified status code and a body
// that is read from the specified file.
func testClientFile(code int, filename string) *Client {
	return &Client{
		http: http.Client{
			Transport: newFileRoundTripper(code, filename),
		},
	}
}

func TestNewReleasesNoAuth(t *testing.T) {
	c := testClientString(400, "")
	_, _, err := c.NewReleases()
	if err == nil {
		t.Errorf("Call should have failed without authorization")
	}
}
