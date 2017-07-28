package spotify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"
	"time"
)

func TestFeaturedPlaylists(t *testing.T) {
	client, server := testClientFile(http.StatusOK, "test_data/featured_playlists.txt")
	defer server.Close()

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
	}
	if p.Playlists == nil || len(p.Playlists) == 0 {
		t.Fatal("Empty playlists result")
	}
	expected := "Hangover Friendly Singer-Songwriter"
	if name := p.Playlists[0].Name; name != expected {
		t.Errorf("Want '%s', got '%s'\n", expected, name)
	}
}

func TestFeaturedPlaylistsExpiredToken(t *testing.T) {
	json := `{
		"error": {
			"status": 401,
			"message": "The access token expired"
		}
	}`
	client, server := testClientString(http.StatusUnauthorized, json)
	defer server.Close()

	msg, pl, err := client.FeaturedPlaylists()
	if msg != "" || pl != nil || err == nil {
		t.Fatal("Expected an error")
	}
	serr, ok := err.(Error)
	if !ok {
		t.Fatalf("Expected spotify Error, got %T", err)
	}
	if serr.Status != http.StatusUnauthorized {
		t.Error("Expected HTTP 401")
	}
}

func TestPlaylistsForUser(t *testing.T) {
	client, server := testClientFile(http.StatusOK, "test_data/playlists_for_user.txt")
	defer server.Close()

	playlists, err := client.GetPlaylistsForUser("whizler")
	if err != nil {
		t.Error(err)
	}
	if l := len(playlists.Playlists); l == 0 {
		t.Fatal("Didn't get any results")
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
	client, server := testClientFile(http.StatusOK, "test_data/get_playlist_opt.txt")
	defer server.Close()

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

func TestFollowPlaylistSetsContentType(t *testing.T) {
	client, server := testClientString(http.StatusOK, "", func(req *http.Request) {
		if req.Header.Get("Content-Type") != "application/json" {
			t.Error("Follow playlist request didn't contain Content-Type: application/json")
		}
	})
	defer server.Close()

	err := client.FollowPlaylist("ownerID", "playlistID", true)
	if err != nil {
		t.Error(err)
	}
}

func TestGetPlaylistTracks(t *testing.T) {
	client, server := testClientFile(http.StatusOK, "test_data/playlist_tracks.txt")
	defer server.Close()

	tracks, err := client.GetPlaylistTracks("user", "playlistID")
	if err != nil {
		t.Error(err)
	}
	if tracks.Total != 47 {
		t.Errorf("Got %d tracks, expected 47\n", tracks.Total)
	}
	if len(tracks.Tracks) == 0 {
		t.Fatal("No tracks returned")
	}
	expected := "Time Of Our Lives"
	actual := tracks.Tracks[0].Track.Name
	if expected != actual {
		t.Errorf("Got '%s', expected '%s'\n", actual, expected)
	}
	added := tracks.Tracks[0].AddedAt
	tm, err := time.Parse(TimestampLayout, added)
	if err != nil {
		t.Error(err)
	}
	if f := tm.Format(DateLayout); f != "2014-11-25" {
		t.Errorf("Expected added at 2014-11-25, got %s\n", f)
	}
}

func TestUserFollowsPlaylist(t *testing.T) {
	client, server := testClientString(http.StatusOK, `[ true, false ]`)
	defer server.Close()

	follows, err := client.UserFollowsPlaylist("jmperezperez", ID("2v3iNvBS8Ay1Gt2uXtUKUT"), "possan", "elogain")
	if err != nil {
		t.Error(err)
	}
	if len(follows) != 2 || !follows[0] || follows[1] {
		t.Errorf("Expected '[true, false]', got %#v\n", follows)
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
	"url": "spotify:user:thelinmichael"
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
"url": "spotify:user:thelinmichael:playlist:7d2D2S200NyUE5KYs80PwO"
}`

func TestCreatePlaylist(t *testing.T) {
	client, server := testClientString(http.StatusCreated, newPlaylist)
	defer server.Close()

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

func TestRenamePlaylist(t *testing.T) {
	client, server := testClientString(http.StatusOK, "")
	defer server.Close()

	if err := client.ChangePlaylistName("user", ID("playlist-id"), "new name"); err != nil {
		t.Error(err)
	}
}

func TestChangePlaylistAccess(t *testing.T) {
	client, server := testClientString(http.StatusOK, "")
	defer server.Close()

	if err := client.ChangePlaylistAccess("user", ID("playlist-id"), true); err != nil {
		t.Error(err)
	}
}

func TestChangePlaylistNamdAndAccess(t *testing.T) {
	client, server := testClientString(http.StatusOK, "")
	defer server.Close()

	if err := client.ChangePlaylistNameAndAccess("user", ID("playlist-id"), "new_name", true); err != nil {
		t.Error(err)
	}
}

func TestChangePlaylistNameFailure(t *testing.T) {
	client, server := testClientString(http.StatusForbidden, "")
	defer server.Close()

	if err := client.ChangePlaylistName("user", ID("playlist-id"), "new_name"); err == nil {
		t.Error("Expected error but didn't get one")
	}
}

func TestAddTracksToPlaylist(t *testing.T) {
	client, server := testClientString(http.StatusCreated, `{ "snapshot_id" : "JbtmHBDBAYu3/bt8BOXKjzKx3i0b6LCa/wVjyl6qQ2Yf6nFXkbmzuEa+ZI/U1yF+" }`)
	defer server.Close()

	snapshot, err := client.AddTracksToPlaylist("user", ID("playlist_id"), ID("track1"), ID("track2"))
	if err != nil {
		t.Error(err)
	}
	if snapshot != "JbtmHBDBAYu3/bt8BOXKjzKx3i0b6LCa/wVjyl6qQ2Yf6nFXkbmzuEa+ZI/U1yF+" {
		t.Error("Didn't get expected snapshot ID")
	}
}

func TestRemoveTracksFromPlaylist(t *testing.T) {
	client, server := testClientString(http.StatusOK, `{ "snapshot_id" : "JbtmHBDBAYu3/bt8BOXKjzKx3i0b6LCa/wVjyl6qQ2Yf6nFXkbmzuEa+ZI/U1yF+" }`, func(req *http.Request) {
		requestBody, err := ioutil.ReadAll(req.Body)

		var body map[string]interface{}
		err = json.Unmarshal(requestBody, &body)
		if err != nil {
			t.Fatal("Error decoding request body:", err)
		}
		tracksArray, ok := body["tracks"]
		if !ok {
			t.Error("No tracks JSON object in request body")
		}
		tracksSlice := tracksArray.([]interface{})
		if l := len(tracksSlice); l != 2 {
			t.Fatalf("Expected 2 tracks, got %d\n", l)
		}
		track0 := tracksSlice[0].(map[string]interface{})
		trackURI, ok := track0["uri"]
		if !ok {
			t.Error("Track object doesn't contain 'uri' field")
		}
		if trackURI != "spotify:track:track1" {
			t.Errorf("Expeced URI: 'spotify:track:track1', got '%s'\n", trackURI)
		}
	})
	defer server.Close()

	snapshotID, err := client.RemoveTracksFromPlaylist("userID", "playlistID", "track1", "track2")
	if err != nil {
		t.Error(err)
	}
	if snapshotID != "JbtmHBDBAYu3/bt8BOXKjzKx3i0b6LCa/wVjyl6qQ2Yf6nFXkbmzuEa+ZI/U1yF+" {
		t.Error("Incorrect snapshot ID")
	}
}

func TestRemoveTracksFromPlaylistOpt(t *testing.T) {
	client, server := testClientString(http.StatusOK, `{ "snapshot_id" : "JbtmHBDBAYu3/bt8BOXKjzKx3i0b6LCa/wVjyl6qQ2Yf6nFXkbmzuEa+ZI/U1yF+" }`, func(req *http.Request) {
		requestBody, err := ioutil.ReadAll(req.Body)
		if err != nil {
			t.Fatal(err)
		}

		var body map[string]interface{}
		err = json.Unmarshal(requestBody, &body)
		if err != nil {
			t.Fatal(err)
		}
		if _, ok := body["snapshot_id"]; ok {
			t.Error("JSON contains snapshot_id field when none was specified")
			fmt.Println(string(requestBody))
			return
		}
		jsonTracks := body["tracks"].([]interface{})
		if len(jsonTracks) != 3 {
			t.Fatal("Expected 3 tracks, got", len(jsonTracks))
		}
		track1 := jsonTracks[1].(map[string]interface{})
		expected := "spotify:track:track1"
		if track1["uri"] != expected {
			t.Fatalf("Want '%s', got '%s'\n", expected, track1["uri"])
		}
		indices := track1["positions"].([]interface{})
		if len(indices) != 1 || int(indices[0].(float64)) != 9 {
			t.Error("Track indices incorrect")
		}
	})
	defer server.Close()

	tracks := []TrackToRemove{
		NewTrackToRemove("track0", []int{0, 4}), // remove track0 in position 0 and 4
		NewTrackToRemove("track1", []int{9}),    // remove track1 in position 9...
		NewTrackToRemove("track2", []int{8}),
	}
	// intentionally not passing a snapshot ID here
	snapshotID, err := client.RemoveTracksFromPlaylistOpt("userID", "playlistID", tracks, "")
	if err != nil || snapshotID != "JbtmHBDBAYu3/bt8BOXKjzKx3i0b6LCa/wVjyl6qQ2Yf6nFXkbmzuEa+ZI/U1yF+" {
		t.Fatal("Remove call failed. err=", err)
	}
}

func TestReplacePlaylistTracks(t *testing.T) {
	client, server := testClientString(http.StatusCreated, "")
	defer server.Close()

	err := client.ReplacePlaylistTracks("userID", "playlistID", "track1", "track2")
	if err != nil {
		t.Error(err)
	}
}

func TestReplacePlaylistTracksForbidden(t *testing.T) {
	client, server := testClientString(http.StatusForbidden, "")
	defer server.Close()

	err := client.ReplacePlaylistTracks("userID", "playlistID", "track1", "track2")
	if err == nil {
		t.Error("Replace succeeded but shouldn't have")
	}
}

func TestReorderPlaylistRequest(t *testing.T) {
	client, server := testClientString(http.StatusNotFound, "", func(req *http.Request) {
		if ct := req.Header.Get("Content-Type"); ct != "application/json" {
			t.Errorf("Expected Content-Type: application/json, got '%s'\n", ct)
		}
		if req.Method != "PUT" {
			t.Errorf("Expected a PUT, got a %s\n", req.Method)
		}
		// unmarshal the JSON into a map[string]interface{}
		// so we can test for existence of certain keys
		var body map[string]interface{}
		json.NewDecoder(req.Body).Decode(&body)

		if start, ok := body["range_start"]; ok {
			if start != float64(3) {
				t.Errorf("Expected range_start to be 3, but it was %#v\n", start)
			}
		} else {
			t.Errorf("Required field range_start is missing")
		}

		if ib, ok := body["insert_before"]; ok {
			if ib != float64(8) {
				t.Errorf("Expected insert_before to be 8, but it was %#v\n", ib)
			}
		} else {
			t.Errorf("Required field insert_before is missing")
		}

		if _, ok := body["range_length"]; ok {
			t.Error("Parameter range_length shouldn't have been in body")
		}
		if _, ok := body["snapshot_id"]; ok {
			t.Error("Parameter snapshot_id shouldn't have been in body")
		}
	})
	defer server.Close()

	client.ReorderPlaylistTracks("user", "playlist", PlaylistReorderOptions{
		RangeStart:   3,
		InsertBefore: 8,
	})
}

func TestSetPlaylistImage(t *testing.T) {
	client, server := testClientString(http.StatusAccepted, "", func(req *http.Request) {
		if ct := req.Header.Get("Content-Type"); ct != "image/jpeg" {
			t.Errorf("wrong content type, got %s, want image/jpeg", ct)
		}
		if req.Method != "PUT" {
			t.Errorf("expected a PUT, got a %s\n", req.Method)
		}
		body, err := ioutil.ReadAll(req.Body)
		if err != nil {
			t.Fatal(err)
		}

		if !bytes.Equal(body, []byte("Zm9v")) {
			t.Errorf("invalid request body: want Zm9v, got %s", string(body))
		}
	})
	defer server.Close()

	err := client.SetPlaylistImage("user", "playlist", bytes.NewReader([]byte("foo")))
	if err != nil {
		t.Fatal(err)
	}
}
