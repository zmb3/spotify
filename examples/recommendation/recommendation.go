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

	// store artist and genre info (upto 5 seed values allowed)
	favorite_artists := []string{"Invent Animate", "Spiritbox", "Sleep Token"}
	genres := []string{"metalcore"}
	var artist_ids []spotify.ID

	// find the ID for each artist and append it to list of artist ID's
	for _, artist := range favorite_artists {
		artists, err := client.Search(ctx, artist, spotify.SearchTypeArtist)

		if err != nil {
			log.Fatal(err)
		}
		artist_ids = append(artist_ids, artists.Artists.Artists[0].ID)
	}

	// store the values in seed
	seeds := spotify.Seeds{
		Artists: artist_ids,
		Genres:  genres,
	}

	// declare track attributes for a more refined search
	track_attributes := spotify.NewTrackAttributes().
		MaxValence(0.4).
		TargetEnergy(0.6).
		TargetDanceability(0.6)

	// get recommendations based on seed values
	res, err := client.GetRecommendations(ctx, seeds, track_attributes, spotify.Country("US"), spotify.Limit(10))
	if err != nil {
		log.Fatal(err)
	}

	// display the recommended tracks along with artists
	fmt.Println("\t---- Recommended Tracks ----")
	for _, track := range res.Tracks {
		fmt.Println(track.Name + " by " + track.Artists[0].Name)
	}
}
