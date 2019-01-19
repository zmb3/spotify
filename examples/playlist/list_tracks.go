package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"golang.org/x/oauth2/clientcredentials"

	"github.com/zmb3/spotify"
)

func main() {
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

	// List names of tracks of that playlist in a single call.
	var SwingdancePlaylistID spotify.ID = "4IqkT7Dviavpz1mF6PKGnA"
	tracks, err := client.GetPlaylistTracksAll(SwingdancePlaylistID)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%d Tracks returned:\n", len(tracks))
	for _, t := range tracks {
		fmt.Printf("\t- %s\n", t.Track.Name)
	}

	fmt.Printf("\n\n---------\n\n")

	// Using client.GetPlaylistTracks will only return the first page of tracks.
	// Page return types return maximum 100 items.
	// They should have a convenience method that returns items from all pages.
	trackpage, err := client.GetPlaylistTracks(SwingdancePlaylistID)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%d Tracks returned from a getting a page:\n", len(trackpage.Tracks))
	for _, t := range trackpage.Tracks {
		fmt.Printf("\t- %s\n", t.Track.Name)
	}
}
