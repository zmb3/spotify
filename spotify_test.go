package spotify

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"testing"
)

func testClient(code int, body string) (*httptest.Server, *Client) {
	baseAddress = "http://localhost/"
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(code)
		fmt.Fprintln(w, body)
	}))
	transport := &http.Transport{
		Proxy: func(req *http.Request) (*url.URL, error) {
			return url.Parse(server.URL)
		},
	}
	client := Client{
		http: http.Client{Transport: transport},
	}
	return server, &client
}

func testClientFromFile(code int, filename string, t *testing.T) (*httptest.Server, *Client) {
	baseAddress = "http://localhost/"
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(code)
		file, err := os.Open(filename)
		if err != nil {
			t.Error(err.Error())
			return
		}
		defer file.Close()
		w.Header().Set("Content-Type", "application/json")
		io.Copy(w, file)
	}))
	transport := &http.Transport{
		Proxy: func(req *http.Request) (*url.URL, error) {
			return url.Parse(server.URL)
		},
	}
	client := Client{
		http: http.Client{Transport: transport},
	}
	return server, &client
}

// func TestSearchNoQuery(t *testing.T) {
// 	client := &Client{}
// 	client.Search("", Artist|Album)
// TODO: expect error 400, no search query
// }

func TestSearchArtist(t *testing.T) {
	server, client := testClientFromFile(http.StatusOK, "test_data/search_artist.txt", t)
	defer server.Close()
	if t.Failed() {
		return
	}
	result, err := client.Search("tania bowra", Artist)
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
	server, client := testClientFromFile(http.StatusOK, "test_data/search_tracks.txt", t)
	defer server.Close()
	if t.Failed() {
		return
	}
	result, err := client.Search("uptown", Track)
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
	server, client := testClientFromFile(http.StatusOK, "test_data/search_trackplaylist.txt", t)
	defer server.Close()
	if t.Failed() {
		return
	}
	result, err := client.Search("holiday", Playlist|Track)
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
