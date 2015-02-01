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

func TestFeaturedPlaylistNoAuth(t *testing.T) {
	var client Client
	_, _, err := client.FeaturedPlaylists()
	if err == nil {
		t.Error("Call should have failed without authorization:", err)
	}
}

func TestFeaturedPlaylists(t *testing.T) {
	client := testClientFile(http.StatusOK, "test_data/featured_playlists.txt")
	addDummyAuth(client)

	country := "SE"
	opt := PlaylistOptions{}
	opt.Country = &country

	msg, p, err := client.FeaturedPlaylistsOpt(&opt)
	if err != nil {
		t.Error(err)
		return
	}
	if msg != "Enjoy a mellow afternoon." {
		t.Errorf("Want 'Enjoy a mellow afternoon.', got'%s'\n", msg)
		return
	}
	if p.Playlists == nil || len(p.Playlists) == 0 {
		t.Error("Empty playlists result")
		return
	}
	expected := "Hangover Friendly Singer-Songwriter"
	if name := p.Playlists[0].Name; name != expected {
		t.Errorf("Want '%s', got '%s'\n", name, expected)
	}
}

func TestFeaturedPlaylistsExpiredToken(t *testing.T) {
	json := `{
		"error": {
			"status": 401,
			"message": "The access token expired"
		}
	}`
	client := testClientString(http.StatusUnauthorized, json)
	addDummyAuth(client)

	msg, pl, err := client.FeaturedPlaylists()
	if msg != "" || pl != nil || err == nil {
		t.Error("Expected an error")
		return
	}
	serr, ok := err.(Error)
	if !ok {
		t.Error("Expected spotify Error")
		return
	}
	if serr.Status != http.StatusUnauthorized {
		t.Error("Expected HTTP 401")
	}
}

func TestPlaylistsForUser(t *testing.T) {
	client := testClientFile(http.StatusOK, "test_data/playlists_for_user.txt")
	addDummyAuth(client)
	playlists, err := client.PlaylistsForUser("whizler")
	if err != nil {
		t.Error(err)
	}
	if l := len(playlists.Playlists); l == 0 {
		t.Error("Didn't get any results")
		return
	}
	p := playlists.Playlists[0]
	if p.Name != "Nederlandse Tipparade" {
		t.Error("Expected Nederlandse Tipparade, got", p.Name)
	}
	if p.Tracks.Total != 29 {
		t.Error("Expected 29 tracks, got", p.Tracks.Total)
	}
}

func TestGetPlaylistOpt(t *testing.T) {
	client := testClientFile(http.StatusOK, "test_data/get_playlist_opt.txt")
	addDummyAuth(client)
	fields := "href,name,owner(!href,external_urls),tracks.items(added_by.id,track(name,href,album(name,href)))"
	p, err := client.GetPlaylistOpt("spotify", "59ZbFPES4DQwEjBpWHzrtC", fields)
	if err != nil {
		t.Error(err)
	}
	if p.Collaborative {
		t.Error("Playlist shouldn't be collaborative")
	}
	if p.Description != "" {
		t.Error("No description should be included")
	}
	if p.Tracks.Total != 10 {
		t.Error("Expected 10 tracks")
	}
}

func TestGetPlaylistTracks(t *testing.T) {
	client := testClientFile(http.StatusOK, "test_data/playlist_tracks.txt")
	addDummyAuth(client)
	tracks, err := client.GetPlaylistTracks("user", "playlistID")
	if err != nil {
		t.Error(err)
	}
	if tracks.Total != 47 {
		t.Errorf("Got %d tracks, expected 47\n", tracks.Total)
	}
	if len(tracks.Tracks) == 0 {
		t.Error("No tracks returned")
		return
	}
	expected := "Time Of Our Lives"
	actual := tracks.Tracks[0].Track.Name
	if expected != actual {
		t.Errorf("Got '%s', expected '%s'\n", actual, expected)
	}
}

var newPlaylist = `
{
"collaborative": false,
"description": null,
"external_urls": {
	"spotify": "http://open.spotify.com/user/thelinmichael/playlist/7d2D2S200NyUE5KYs80PwO"
},
"followers": {
	"href": null,
	"total": 0
},
"href": "https://api.spotify.com/v1/users/thelinmichael/playlists/7d2D2S200NyUE5KYs80PwO",
"id": "7d2D2S200NyUE5KYs80PwO",
"images": [ ],
"name": "A New Playlist",
"owner": {
	"external_urls": {
	"spotify": "http://open.spotify.com/user/thelinmichael"
	},
	"href": "https://api.spotify.com/v1/users/thelinmichael",
	"id": "thelinmichael",
	"type": "user",
	"uri": "spotify:user:thelinmichael"
},
"public": false,
"snapshot_id": "s0o3TSuYnRLl2jch+oA4OEbKwq/fNxhGBkSPnvhZdmWjNV0q3uCAWuGIhEx8SHIx",
"tracks": {
	"href": "https://api.spotify.com/v1/users/thelinmichael/playlists/7d2D2S200NyUE5KYs80PwO/tracks",
	"items": [ ],
	"limit": 100,
	"next": null,
	"offset": 0,
	"previous": null,
	"total": 0
},
"type": "playlist",
"uri": "spotify:user:thelinmichael:playlist:7d2D2S200NyUE5KYs80PwO"
}`

func TestCreatePlaylist(t *testing.T) {
	client := testClientString(http.StatusCreated, newPlaylist)
	addDummyAuth(client)
	p, err := client.CreatePlaylistForUser("thelinmichael", "A New Playlist", false)
	if err != nil {
		t.Error(err)
	}
	if p.IsPublic {
		t.Error("Expected private playlist, got public")
	}
	if p.Name != "A New Playlist" {
		t.Errorf("Expected 'A New Playlist', got '%s'\n", p.Name)
	}
	if p.Tracks.Total != 0 {
		t.Error("Expected new playlist to be empty")
	}
}
