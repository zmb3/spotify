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
	// --- (1) ----
	// Login to Spotify
	ctx := context.Background()

	config := &clientcredentials.Config{
		ClientID:     os.Getenv("SPOTIFY_ID"),
		ClientSecret: os.Getenv("SPOTIFY_SECRET"),
		TokenURL:     spotifyauth.TokenURL,
		Scopes:       []string{spotifyauth.ScopeUserReadPrivate, spotifyauth.ScopeUserReadEmail},
	}
	token, err := config.Token(ctx)
	if err != nil {
		log.Fatalf("couldn't get token: %v", err)
	}

	httpClient := spotifyauth.Authenticator{}.Client(ctx, token)
	client := spotify.New(httpClient)

	// --- (2) ----
	// Run a search and recover all the playlists that have christmas in their name
	results, err := client.Search(ctx, "christmas", spotify.SearchTypePlaylist)
	if err != nil {
		log.Fatal(err)
	}

	// --- (3) ----
	// Handle playlist results
	tracks := make(map[string]spotify.PlaylistTrack, 40)
	if results.Playlists != nil {
		for _, item := range results.Playlists.Playlists {

			// --- (3.1) ----
			// Get all the songs from the playlists
			playlistsTracks, err := client.GetPlaylistTracks(ctx, item.ID)
			if err != nil {
				log.Fatal(err)
			}
			for page := 1; ; page++ {
				for _, track := range playlistsTracks.Tracks {
					// If the track is not already included, we add it
					if _, found := tracks[track.Track.Name]; !found && len(tracks) < 40 {
						tracks[track.Track.Name] = track
					}
				}
				err = client.NextPage(ctx, playlistsTracks)
				if err == spotify.ErrNoMorePages {
					break
				}
				if err != nil {
					log.Fatal(err)
				}
			}
		}
	}

	fmt.Println("Top Christams Songs:")
	for k := range tracks {
		fmt.Println("   ", k)
	}

}
