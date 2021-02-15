package main

import (
	"context"
	"log"
	"os"

	"github.com/zmb3/spotify"
	"golang.org/x/oauth2/clientcredentials"
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
	tracks, err := client.GetPlaylistTracks("57qttz6pK881sjxj2TAEEo")
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Playlist has %d total tracks", tracks.Total)
	for page := 1; ; page++ {
		log.Printf("  Page %d has %d tracks", page, len(tracks.Tracks))
		err = client.NextPage(tracks)
		if err == spotify.ErrNoMorePages {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
	}
}
