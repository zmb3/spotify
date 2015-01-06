package spotify

import "testing"

// func TestSearchNoQuery(t *testing.T) {
// 	client := &Client{}
// 	client.Search("", Artist|Album)
// TODO: expect error 400, no search query
// }

func TestSearchArtist(t *testing.T) {
	client := &Client{}
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
	if result.Artists == nil {
		t.Error("Didn't receive artist results")
	}
}

func TestSearchTracks(t *testing.T) {
	client := &Client{}
	result, err := client.Search("uptown", Track)
	if err != nil {
		t.Error(err)
	}
	if result.Tracks == nil {
		t.Error("Didn't receive track results")
	}
}

func TestSearchAlbumPlaylist(t *testing.T) {
	client := &Client{}
	result, err := client.Search("holiday", Playlist|Album)
	if err != nil {
		t.Error(err)
	}
	if result.Tracks != nil {
		t.Error("Searched for albums and playlists but received track results")
	}
	if result.Artists != nil {
		t.Error("Searched for albums and playlists but received artist results")
	}
	if result.Albums == nil {
		t.Error("Didn't receive album results")
	}
	if result.Playlists == nil {
		t.Error("Didn't receive playlist results")
	}
}
