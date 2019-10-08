package main

import (
	"testing"

	"github.com/zmb3/spotify"
)

type mockSpotifyClient struct{}

func (m *mockSpotifyClient) GetPlaylist(playlistID spotify.ID) (*spotify.FullPlaylist, error) {
	return &spotify.FullPlaylist{
		SimplePlaylist: spotify.SimplePlaylist{
			Name: "whatever",
		},
	}, nil
}

func Test_NewGetPlaylistName(t *testing.T) {
	client := &mockSpotifyClient{}

	name := getPlaylistName(client, "whatever")

	if name != "whatever" {
		t.Errorf("expected %s, got %s", "whatever", name)
	}
}
