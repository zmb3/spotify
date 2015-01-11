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

// func TestSearchNoQuery(t *testing.T) {
// 	client := &Client{}
// 	client.Search("", Artist|Album)
// TODO: expect error 400, no search query
// }

func TestSearchArtist(t *testing.T) {
	client := testClientFile(http.StatusOK, "test_data/search_artist.txt")
	result, err := client.Search("tania bowra", SearchTypeArtist)
	if err != nil {
		t.Error(err)
	}
	if result.Albums != nil {
		t.Error("Searched for artists but received album results")
	}
	if result.Playlists != nil {
		t.Error("Searched for artists but received playlist results")
	}
	if result.Tracks != nil {
		t.Error("Searched for artists but received track results")
	}
	if result.Artists == nil || len(result.Artists.Artists) == 0 {
		t.Error("Didn't receive artist results")
	}
	if result.Artists.Artists[0].Name != "Tania Bowra" {
		t.Error("Got wrong artist name")
	}
}

func TestSearchTracks(t *testing.T) {
	client := testClientFile(http.StatusOK, "test_data/search_tracks.txt")
	result, err := client.Search("uptown", SearchTypeTrack)
	if err != nil {
		t.Error(err)
	}
	if result.Tracks == nil || len(result.Tracks.Tracks) == 0 {
		t.Error("Didn't receive track results")
	}
	if result.Albums != nil {
		t.Error("Searched for tracks but got album results")
	}
	if result.Playlists != nil {
		t.Error("Searched for tracks but got playlist results")
	}
	if result.Artists != nil {
		t.Error("Searched for tracks but got artist results")
	}
	if name := result.Tracks.Tracks[0].Name; name != "Uptown Funk" {
		t.Errorf("Got %s, wanted Uptown Funk\n", name)
	}
}

func TestSearchPlaylistTrack(t *testing.T) {
	client := testClientFile(http.StatusOK, "test_data/search_trackplaylist.txt")
	result, err := client.Search("holiday", SearchTypePlaylist|SearchTypeTrack)
	if err != nil {
		t.Error(err)
	}
	if result.Albums != nil {
		t.Error("Searched for playlists and tracks but received album results")
	}
	if result.Artists != nil {
		t.Error("Searched for playlists and tracks but received artist results")
	}
	if result.Tracks == nil {
		t.Error("Didn't receive track results")
	}
	if result.Playlists == nil {
		t.Error("Didn't receive playlist results")
	}
}
