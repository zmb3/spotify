package spotify

import (
	"net/http"
	"testing"
)

// The example from https://developer.spotify.com/web-api/get-album/
func TestFindAlbum(t *testing.T) {
	client := testClientFile(http.StatusOK, "test_data/find_album.txt")
	album, err := client.GetAlbum(ID("0sNOF9WDwhWunNAHPD3Baj"))
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
	released := album.ReleaseDateTime()
	if released.Year() != 1983 {
		t.Errorf("Expected release date 1983, got %d\n", released.Year())
	}
}

func TestFindAlbumBadID(t *testing.T) {
	client := testClientString(http.StatusNotFound, `{ "error": { "status": 404, "message": "non existing id" } }`)

	album, err := client.GetAlbum(ID("asdf"))
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
	client := testClientFile(http.StatusOK, "test_data/find_albums.txt")
	res, err := client.GetAlbums(ID("41MnTivkwTO3UUJ8DrqEJJ"), ID("6JWc4iAiJ9FjyK0B59ABb4"), ID("6UXCm6bOO4gFlDQZV5yL37"), ID("0X8vBD8h1Ga9eLT8jx9VCC"))
	if err != nil {
		t.Error(err.Error())
		return
	}
	if len(res) != 4 {
		t.Errorf("Expected 4 albums, got %d", len(res))
		return
	}
	expectedAlbums := []string{
		"The Best Of Keane (Deluxe Edition)",
		"Strangeland",
		"Night Train",
		"Mirrored",
	}
	for i, name := range expectedAlbums {
		if res[i].Name != name {
			t.Error("Expected album", name, "but got", res[i].Name)
		}
	}
	release := res[0].ReleaseDateTime()
	if release.Year() != 2013 ||
		release.Month() != 11 ||
		release.Day() != 8 {
		t.Errorf("Expected release 2013-11-08, got %d-%02d-%02d\n",
			release.Year(), release.Month(), release.Day())
	}
	releaseMonthPrecision := res[3].ReleaseDateTime()
	if releaseMonthPrecision.Year() != 2007 ||
		releaseMonthPrecision.Month() != 3 ||
		releaseMonthPrecision.Day() != 1 {
		t.Errorf("Expected release 2007-03-01, got %d-%02d-%02d\n",
			releaseMonthPrecision.Year(), releaseMonthPrecision.Month(), releaseMonthPrecision.Day())
	}
}

func TestFindAlbumTracks(t *testing.T) {
	client := testClientFile(http.StatusOK, "test_data/find_album_tracks.txt")
	res, err := client.GetAlbumTracksOpt(ID("0sNOF9WDwhWunNAHPD3Baj"), 1, 0)
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
