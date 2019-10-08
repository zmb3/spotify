// This example shows how to mock the spotify.Client for unit tests using interfaces
// Full explanation in this blog post: https://deployeveryday.com/2019/10/08/golang-auth-mock.html

package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/zmb3/spotify"
	"golang.org/x/oauth2/clientcredentials"
)

type spotifyClient interface {
	GetPlaylist(playlistID spotify.ID) (*spotify.FullPlaylist, error)
}

func main() {
	const PLAYLIST_ID spotify.ID = "4OyKDT6cLw96G7bd8nTfxD"

	client := newSpotifyClient()
	name := getPlaylistName(client, PLAYLIST_ID)

	fmt.Println(name)
}

func newSpotifyClient() *spotify.Client {
	config := &clientcredentials.Config{
		ClientID:     os.Getenv("SPOTIFY_ID"),
		ClientSecret: os.Getenv("SPOTIFY_SECRET"),
		TokenURL:     spotify.TokenURL,
	}
	token, err := config.Token(context.Background())
	if err != nil {
		log.Fatalf("couldn't get token: %v", err)
	}

	client := spotify.Authenticator{}.NewClient(token)

	return &client
}

func getPlaylistName(client spotifyClient, playlistID spotify.ID) string {
	result, err := client.GetPlaylist(playlistID)
	if err != nil {
		log.Fatalf("couldn't get playlist: %v", err)
	}

	return result.Name
}
