// This example demonstrates how to authenticate with Spotify.
// In order to run this example yourself, you'll need to:
//
//  1. Register an application at: https://developer.spotify.com/my-applications/
//       - Use "http://localhost:8080/callback" as the redirect URI
//  2. Set the SPOTIFY_ID environment variable to the client ID you got in step 1.
//  3. Set the SPOTIFY_SECRET environment variable to the client secret from step 1.
package main

import (
	"context"
	"fmt"
	spotifyauth "github.com/zmb3/spotify/v2/auth"
	"log"
	"net/http"
	"strings"

	"github.com/zmb3/spotify/v2"
)

// redirectURI is the OAuth redirect URI for the application.
// You must register an application at Spotify's developer portal
// and enter this value.
const redirectURI = "http://localhost:8080/callback"

var html = `
<br/>
<a href="/player/play">Play</a><br/>
<a href="/player/pause">Pause</a><br/>
<a href="/player/next">Next track</a><br/>
<a href="/player/previous">Previous Track</a><br/>
<a href="/player/shuffle">Shuffle</a><br/>

`

var (
	auth  = spotifyauth.New(
		spotifyauth.WithRedirectURL(redirectURI),
		spotifyauth.WithScopes(spotifyauth.ScopeUserReadCurrentlyPlaying, spotifyauth.ScopeUserReadPlaybackState, spotifyauth.ScopeUserModifyPlaybackState),
		)
	ch    = make(chan *spotify.Client)
	state = "abc123"
)

func main() {
	// We'll want these variables sooner rather than later
	var client *spotify.Client
	var playerState *spotify.PlayerState

	http.HandleFunc("/callback", completeAuth)

	http.HandleFunc("/player/", func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		action := strings.TrimPrefix(r.URL.Path, "/player/")
		fmt.Println("Got request for:", action)
		var err error
		switch action {
		case "play":
			err = client.Play(ctx)
		case "pause":
			err = client.Pause(ctx)
		case "next":
			err = client.Next(ctx)
		case "previous":
			err = client.Previous(ctx)
		case "shuffle":
			playerState.ShuffleState = !playerState.ShuffleState
			err = client.Shuffle(ctx, playerState.ShuffleState)
		}
		if err != nil {
			log.Print(err)
		}

		w.Header().Set("Content-Type", "text/html")
		fmt.Fprint(w, html)
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Println("Got request for:", r.URL.String())
	})

	go func() {
		url := auth.AuthURL(state)
		fmt.Println("Please log in to Spotify by visiting the following page in your browser:", url)

		// wait for auth to complete
		client = <-ch

		// use the client to make calls that require authorization
		user, err := client.CurrentUser(context.Background())
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("You are logged in as:", user.ID)

		playerState, err = client.PlayerState(context.Background())
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Found your %s (%s)\n", playerState.Device.Type, playerState.Device.Name)
	}()

	http.ListenAndServe(":8080", nil)

}

func completeAuth(w http.ResponseWriter, r *http.Request) {
	tok, err := auth.Token(r.Context(), state, r.URL)
	if err != nil {
		http.Error(w, "Couldn't get token", http.StatusForbidden)
		log.Fatal(err)
	}
	if st := r.FormValue("state"); st != state {
		http.NotFound(w, r)
		log.Fatalf("State mismatch: %s != %s\n", st, state)
	}
	// use the token to get an authenticated client
	client := spotify.New(auth.Client(r.Context(), tok))
	w.Header().Set("Content-Type", "text/html")
	fmt.Fprintf(w, "Login Completed!"+html)
	ch <- client
}
