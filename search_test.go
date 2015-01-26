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
