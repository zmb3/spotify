package spotify

import "testing"

// The example from https://developer.spotify.com/web-api/get-album/
func TestFindAlbum(t *testing.T) {
	server, client := testClientFromFile(200, "test_data/find_album.txt", t)
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

// The example from https://developer.spotify.com/web-api/get-several-albums/
func TestFindAlbums(t *testing.T) {
	server, client := testClientFromFile(200, "test_data/find_albums.txt", t)
	defer server.Close()
	if t.Failed() {
		return
	}
	res, err := client.FindAlbums(ID("41MnTivkwTO3UUJ8DrqEJJ"), ID("6JWc4iAiJ9FjyK0B59ABb4"), ID("6UXCm6bOO4gFlDQZV5yL37"))
	if err != nil {
		t.Error(err.Error())
		return
	}
	if len(res.Albums) != 3 {
		t.Error("Expected 3 albums, got " + string(len(res.Albums)))
		return
	}
	expectedAlbums := []string{
		"The Best Of Keane (Deluxe Edition)",
		"Strangeland",
		"Night Train",
	}
	for i, name := range expectedAlbums {
		if res.Albums[i].Name != name {
			t.Error("Expected album", name, "but got", res.Albums[i].Name)
		}
	}
}

func TestFindAlbumTracks(t *testing.T) {
	server, client := testClientFromFile(200, "test_data/find_album_tracks.txt", t)
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
