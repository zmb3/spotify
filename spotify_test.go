package spotify

import (
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

func testClient(code int, body io.Reader, validators ...func(*http.Request)) (*Client, *httptest.Server) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		for _, v := range validators {
			v(r)
		}
		w.WriteHeader(code)
		io.Copy(w, body)
		r.Body.Close()
		if closer, ok := body.(io.Closer); ok {
			closer.Close()
		}
	}))
	client := &Client{
		http:    http.DefaultClient,
		baseURL: server.URL + "/",
	}
	return client, server
}

// Returns a client whose requests will always return
// the specified status code and body.
func testClientString(code int, body string, validators ...func(*http.Request)) (*Client, *httptest.Server) {
	return testClient(code, strings.NewReader(body))
}

// Returns a client whose requests will always return
// a response with the specified status code and a body
// that is read from the specified file.
func testClientFile(code int, filename string, validators ...func(*http.Request)) (*Client, *httptest.Server) {
	f, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	return testClient(code, f)
}

func TestNewReleases(t *testing.T) {
	c, s := testClientFile(http.StatusOK, "test_data/new_releases.txt")
	defer s.Close()

	r, err := c.NewReleases()
	if err != nil {
		t.Fatal(err)
	}
	if r.Albums[0].ID != "60mvULtYiNSRmpVvoa3RE4" {
		t.Error("Invalid data:", r.Albums[0].ID)
	}
	if r.Albums[0].Name != "We Are One (Ole Ola) [The Official 2014 FIFA World Cup Song]" {
		t.Error("Invalid data", r.Albums[0].Name)
	}
}

func TestNewReleasesRateLimitExceeded(t *testing.T) {
	t.Parallel()
	handlers := []http.HandlerFunc{
		// first attempt fails
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Retry-After", "2")
			w.WriteHeader(rateLimitExceededStatusCode)
			io.WriteString(w, `{ "error": { "message": "slow down", "status": 429 } }`)
		}),
		// next attempt succeeds
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			f, err := os.Open("test_data/new_releases.txt")
			if err != nil {
				t.Fatal(err)
			}
			defer f.Close()
			_, err = io.Copy(w, f)
			if err != nil {
				t.Fatal(err)
			}
		}),
	}

	i := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handlers[i](w, r)
		i++
	}))
	defer server.Close()

	client := &Client{http: http.DefaultClient, baseURL: server.URL + "/", AutoRetry: true}
	releases, err := client.NewReleases()
	if err != nil {
		t.Fatal(err)
	}
	if releases.Albums[0].ID != "60mvULtYiNSRmpVvoa3RE4" {
		t.Error("Invalid data:", releases.Albums[0].ID)
	}
}
