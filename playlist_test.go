package spotify

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"
	"time"
)

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
	playlists, err := client.GetPlaylistsForUser("whizler")
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

func TestFollowPlaylistSetsContentType(t *testing.T) {
	client := testClientString(http.StatusOK, "")
	addDummyAuth(client)
	err := client.FollowPlaylist("ownerID", "playlistID", true)
	if err != nil {
		t.Error(err)
		return
	}
	req := getLastRequest(client)
	if req == nil {
		t.Error("Last request was nil")
		return
	}
	if req.Header.Get("Content-Type") != "application/json" {
		t.Error("Follow playlist request didn't contain Content-Type: application/json")
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
	client := testClientString(http.StatusOK, `[ true, false ]`)
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

func TestRenamePlaylist(t *testing.T) {
	client := testClientString(http.StatusOK, "")
	addDummyAuth(client)
	if err := client.ChangePlaylistName("user", ID("playlist-id"), "new name"); err != nil {
		t.Error(err)
	}
}

func TestChangePlaylistAccess(t *testing.T) {
	client := testClientString(http.StatusOK, "")
	addDummyAuth(client)
	if err := client.ChangePlaylistAccess("user", ID("playlist-id"), true); err != nil {
		t.Error(err)
	}
}

func TestChangePlaylistNamdAndAccess(t *testing.T) {
	client := testClientString(http.StatusOK, "")
	addDummyAuth(client)
	if err := client.ChangePlaylistNameAndAccess("user", ID("playlist-id"), "new_name", true); err != nil {
		t.Error(err)
	}
}

func TestChangePlaylistNameFailure(t *testing.T) {
	client := testClientString(http.StatusForbidden, "")
	addDummyAuth(client)
	if err := client.ChangePlaylistName("user", ID("playlist-id"), "new_name"); err == nil {
		t.Error("Expected error but didn't get one")
	}
}

func TestAddTracksToPlaylist(t *testing.T) {
	client := testClientString(http.StatusCreated, `{ "snapshot_id" : "JbtmHBDBAYu3/bt8BOXKjzKx3i0b6LCa/wVjyl6qQ2Yf6nFXkbmzuEa+ZI/U1yF+" }`)
	addDummyAuth(client)
	snapshot, err := client.AddTracksToPlaylist("user", ID("playlist_id"), ID("track1"), ID("track2"))
	if err != nil {
		t.Error(err)
	}
	if snapshot != "JbtmHBDBAYu3/bt8BOXKjzKx3i0b6LCa/wVjyl6qQ2Yf6nFXkbmzuEa+ZI/U1yF+" {
		t.Error("Didn't get expected snapshot ID")
	}
}

func TestRemoveTracksFromPlaylist(t *testing.T) {
	client := testClientString(http.StatusOK, `{ "snapshot_id" : "JbtmHBDBAYu3/bt8BOXKjzKx3i0b6LCa/wVjyl6qQ2Yf6nFXkbmzuEa+ZI/U1yF+" }`)
	addDummyAuth(client)
	snapshotID, err := client.RemoveTracksFromPlaylist("userID", "playlistID", "track1", "track2")
	if err != nil {
		t.Error(err)
	}
	if snapshotID != "JbtmHBDBAYu3/bt8BOXKjzKx3i0b6LCa/wVjyl6qQ2Yf6nFXkbmzuEa+ZI/U1yF+" {
		t.Error("Incorrect snapshot ID")
	}

	req := getLastRequest(client)
	if req == nil {
		t.Error("Last request was nil")
		return
	}
	requestBody, err := ioutil.ReadAll(req.Body)
	req.Body.Close()

	var body map[string]interface{}
	err = json.Unmarshal(requestBody, &body)
	if err != nil {
		t.Error("Error decoding request body:", err)
		return
	}
	tracksArray, ok := body["tracks"]
	if !ok {
		t.Error("No tracks JSON object in request body")
		return
	}
	tracksSlice := tracksArray.([]interface{})
	if l := len(tracksSlice); l != 2 {
		t.Errorf("Expected 2 tracks, got %d\n", l)
		return
	}
	track0 := tracksSlice[0].(map[string]interface{})
	trackURI, ok := track0["uri"]
	if !ok {
		t.Error("Track object doesn't contain 'uri' field")
		return
	}
	if trackURI != "spotify:track:track1" {
		t.Errorf("Expeced URI: 'spotify:track:track1', got '%s'\n", trackURI)
	}
}

func TestRemoveTracksFromPlaylistOpt(t *testing.T) {
	client := testClientString(http.StatusOK, `{ "snapshot_id" : "JbtmHBDBAYu3/bt8BOXKjzKx3i0b6LCa/wVjyl6qQ2Yf6nFXkbmzuEa+ZI/U1yF+" }`)
	addDummyAuth(client)
	tracks := []TrackToRemove{
		NewTrackToRemove("track0", []int{0, 4}), // remove track0 in position 0 and 4
		NewTrackToRemove("track1", []int{9}),    // remove track1 in position 9...
		NewTrackToRemove("track2", []int{8}),
	}
	// intentionally not passing a snapshot ID here
	snapshotID, err := client.RemoveTracksFromPlaylistOpt("userID", "playlistID", tracks, "")
	if err != nil || snapshotID != "JbtmHBDBAYu3/bt8BOXKjzKx3i0b6LCa/wVjyl6qQ2Yf6nFXkbmzuEa+ZI/U1yF+" {
		t.Error("Remove call failed. err=", err)
		return
	}
	// now make sure we got the JSON correct
	req := getLastRequest(client)
	if req == nil {
		t.Error("last request was nil")
		return
	}
	requestBody, err := ioutil.ReadAll(req.Body)
	if err != nil {
		t.Error(err)
		return
	}

	var body map[string]interface{}
	err = json.Unmarshal(requestBody, &body)
	if err != nil {
		t.Error(err)
		return
	}
	if _, ok := body["snapshot_id"]; ok {
		t.Error("JSON contains snapshot_id field when none was specified")
		fmt.Println(string(requestBody))
		return
	}
	jsonTracks := body["tracks"].([]interface{})
	if len(jsonTracks) != 3 {
		t.Error("Expected 3 tracks, got", len(jsonTracks))
		return
	}
	track1 := jsonTracks[1].(map[string]interface{})
	expected := "spotify:track:track1"
	if track1["uri"] != expected {
		t.Errorf("Want '%s', got '%s'\n", expected, track1["uri"])
		return
	}
	indices := track1["positions"].([]interface{})
	if len(indices) != 1 || int(indices[0].(float64)) != 9 {
		t.Error("Track indices incorrect")
	}
}

func TestReplacePlaylistTracks(t *testing.T) {
	client := testClientString(http.StatusCreated, "")
	addDummyAuth(client)
	err := client.ReplacePlaylistTracks("userID", "playlistID", "track1", "track2")
	if err != nil {
		t.Error(err)
	}
}

func TestReplacePlaylistTracksForbidden(t *testing.T) {
	client := testClientString(http.StatusForbidden, "")
	addDummyAuth(client)
	err := client.ReplacePlaylistTracks("userID", "playlistID", "track1", "track2")
	if err == nil {
		t.Error("Replace succeeded but shouldn't have")
	}
}

func TestReorderPlaylistRequest(t *testing.T) {
	client := testClientString(http.StatusNotFound, "")
	userID := "user"
	client.ReorderPlaylistTracks(userID, "playlist", PlaylistReorderOptions{
		RangeStart:   3,
		InsertBefore: 8,
	})
	req := getLastRequest(client)
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
}
