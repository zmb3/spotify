package spotify

import (
	"context"
	"net/http"
	"testing"
)

// The example from https://developer.spotify.com/web-api/get-album/
func TestFindAlbum(t *testing.T) {
	client, server := testClientFile(http.StatusOK, "test_data/find_album.txt")
	defer server.Close()

	album, err := client.GetAlbum(context.Background(), ID("0sNOF9WDwhWunNAHPD3Baj"))
	if err != nil {
		t.Fatal(err)
	}
	if album == nil {
		t.Fatal("Got nil album")
	}
	if album.Name != "She's So Unusual" {
		t.Error("Got wrong album")
	}
	released := album.ReleaseDateTime()
	if released.Year() != 1983 {
		t.Errorf("Expected release date 1983, got %d\n", released.Year())
	}
}

func TestFindAlbumBadID(t *testing.T) {
	client, server := testClientString(http.StatusNotFound, `{ "error": { "status": 404, "message": "non existing id" } }`)
	defer server.Close()

	album, err := client.GetAlbum(context.Background(), ID("asdf"))
	if album != nil {
		t.Fatal("Expected nil album, got", album.Name)
	}
	se, ok := err.(Error)
	if !ok {
		t.Error("Expected spotify error, got", err)
	}
	if se.Status != 404 {
		t.Errorf("Expected HTTP 404, got %d. ", se.Status)
	}
	if se.Message != "non existing id" {
		t.Error("Unexpected error message: ", se.Message)
	}
}

func TestFindAlbumTracks(t *testing.T) {
	client, server := testClientFile(http.StatusOK, "test_data/find_album_tracks.txt")
	defer server.Close()

	res, err := client.GetAlbumTracks(context.Background(), ID("0sNOF9WDwhWunNAHPD3Baj"), Limit(1))
	if err != nil {
		t.Fatal(err)
	}
	if res.Total != 13 {
		t.Fatal("Got", res.Total, "results, want 13")
	}
	if len(res.Tracks) == 1 {
		if res.Tracks[0].Name != "Money Changes Everything" {
			t.Error("Expected track 'Money Changes Everything', got", res.Tracks[0].Name)
		}
	} else {
		t.Error("Expected 1 track, got", len(res.Tracks))
	}
}
