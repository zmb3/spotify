// This example demonstrates how to authenticate with Spotify using the
// client credentials flow.  Note that this flow does not include authorization
// and can't be used to access a user's private data.
//
// Make sure you set the SPOTIFY_ID and SPOTIFY_SECRET environment variables
// prior to running this example.
package main

import (
	"context"
	"fmt"
	"github.com/zmb3/spotify/v2/auth"
	"log"
	"os"

	"github.com/zmb3/spotify/v2"
	"golang.org/x/oauth2/clientcredentials"
)

func main() {
	ctx := context.Background()
	config := &clientcredentials.Config{
		ClientID:     os.Getenv("SPOTIFY_ID"),
		ClientSecret: os.Getenv("SPOTIFY_SECRET"),
		TokenURL:     spotifyauth.TokenURL,
	}
	token, err := config.Token(ctx)
	if err != nil {
		log.Fatalf("couldn't get token: %v", err)
	}

	httpClient := spotifyauth.New("").Client(ctx, token)
	client := spotify.New(spotify.HTTPClientOpt(httpClient))
	msg, page, err := client.FeaturedPlaylists(ctx)
	if err != nil {
		log.Fatalf("couldn't get features playlists: %v", err)
	}

	fmt.Println(msg)
	for _, playlist := range page.Playlists {
		fmt.Println("  ", playlist.Name)
	}
}
