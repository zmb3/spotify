package main

import (
	"context"
	"fmt"
	"log"
	"os"

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
	artists, err := client.Search(ctx, "Invent Animate", spotify.SearchTypeArtist)
	if err != nil {
		log.Fatal(err)
	}

	// select the top album
	artist := artists.Artists.Artists[0]

	seeds := spotify.Seeds{
		Artists: []spotify.ID{artist.ID},
		Genres:  []string{"metalcore", "hardcore"},
	}

	res, err := client.GetRecommendations(ctx, seeds, nil)
	if err != nil {
		log.Fatal(err)
	}
	for _, rs := range res.Seeds {
		fmt.Println(rs)
	}
	for _, st := range res.Tracks {
		fmt.Println(st.Name + " - " + st.Artists[0].Name)
	}
}
