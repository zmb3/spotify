// Copyright 2014, 2015 Zac Bergquist
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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
	res, err := client.GetAlbums(ID("41MnTivkwTO3UUJ8DrqEJJ"), ID("6JWc4iAiJ9FjyK0B59ABb4"), ID("6UXCm6bOO4gFlDQZV5yL37"))
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
