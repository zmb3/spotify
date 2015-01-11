package spotify

import (
	"net/http"
	"testing"
)

// The example from https://developer.spotify.com/web-api/get-album/
func TestFindAlbum(t *testing.T) {
	server, client := testClientFromFile(http.StatusOK, "test_data/find_album.txt", t)
	defer server.Close()
	if t.Failed() {
		return
	}
	album, err := client.FindAlbum(ID("0sNOF9WDwhWunNAHPD3Baj"))
	if err != nil {
		t.Error(err.Error())
		return
	}
	if album == nil {
		t.Error("Got nil album")
		return
	}
	if album.Name != "She's So Unusual" {
		t.Error("Got wrong album")
	}
}

func TestFindAlbumBadID(t *testing.T) {
	server, client := testClient(404, `{ "error": { "status": 404, "message": "non existing id" } }`)
	defer server.Close()

	album, err := client.FindAlbum(ID("asdf"))
	if album != nil {
		t.Error("Expected nil album, got", album.Name)
		return
	}
	se, ok := err.(Error)
	if !ok {
		t.Error("Expected spotify error, got", err)
		return
	}
	if se.Status != 404 {
		t.Errorf("Expected HTTP 404, got %d. ", se.Status)
		return
	}
	if se.Message != "non existing id" {
		t.Error("Unexpected error message: ", se.Message)
	}
}

// The example from https://developer.spotify.com/web-api/get-several-albums/
func TestFindAlbums(t *testing.T) {
	server, client := testClientFromFile(http.StatusOK, "test_data/find_albums.txt", t)
	defer server.Close()
	if t.Failed() {
		return
	}
	res, err := client.FindAlbums(ID("41MnTivkwTO3UUJ8DrqEJJ"), ID("6JWc4iAiJ9FjyK0B59ABb4"), ID("6UXCm6bOO4gFlDQZV5yL37"))
	if err != nil {
		t.Error(err.Error())
		return
	}
	if len(res) != 3 {
		t.Error("Expected 3 albums, got " + string(len(res)))
		return
	}
	expectedAlbums := []string{
		"The Best Of Keane (Deluxe Edition)",
		"Strangeland",
		"Night Train",
	}
	for i, name := range expectedAlbums {
		if res[i].Name != name {
			t.Error("Expected album", name, "but got", res[i].Name)
		}
	}
}

func TestFindAlbumTracks(t *testing.T) {
	server, client := testClientFromFile(http.StatusOK, "test_data/find_album_tracks.txt", t)
	defer server.Close()
	if t.Failed() {
		return
	}
	res, err := client.FindAlbumTracksLimited(ID("0sNOF9WDwhWunNAHPD3Baj"), 1, 0)
	if err != nil {
		t.Error(err)
		return
	}
	if res.Total != 13 {
		t.Error("Got", res.Total, "results, want 13")
	}
	if len(res.Tracks) == 1 {
		if res.Tracks[0].Name != "Money Changes Everything" {
			t.Error("Expected track 'Money Changes Everything', got", res.Tracks[0].Name)
		}
	} else {
		t.Error("Expected 1 track, got", len(res.Tracks))
	}
}
