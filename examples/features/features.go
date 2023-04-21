package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"text/tabwriter"

	spotifyauth "github.com/zmb3/spotify/v2/auth"

	"golang.org/x/oauth2/clientcredentials"

	"github.com/zmb3/spotify/v2"
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

	httpClient := spotifyauth.New().Client(ctx, token)
	client := spotify.New(httpClient)

	// search for albums with the name Sempiternal
	results, err := client.Search(ctx, "Sempiternal", spotify.SearchTypeAlbum)
	if err != nil {
		log.Fatal(err)
	}

	// select the top album
	item := results.Albums.Albums[0]

	// get tracks from album
	res, err := client.GetAlbumTracks(ctx, item.ID, spotify.Market("US"))

	if err != nil {
		log.Fatal("error getting tracks ....", err.Error())
		return
	}

	// *display in tabular form using TabWriter
	w := tabwriter.NewWriter(os.Stdout, 10, 2, 3, ' ', 0)
	fmt.Fprintf(w, "%s\t%s\t%s\t%s\t\n\n", "Songs", "Energy", "Danceability", "Valence")

	// loop through tracks
	for _, track := range res.Tracks {

		// retrieve features
		features, err := client.GetAudioFeatures(ctx, track.ID)
		if err != nil {
			log.Fatal("error getting audio features...", err.Error())
			return
		}
		fmt.Fprintf(w, "%s\t%v\t%v\t%v\t\n", track.Name, features[0].Energy, features[0].Danceability, features[0].Valence)
		w.Flush()
	}
}
