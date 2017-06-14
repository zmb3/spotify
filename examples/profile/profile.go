// Command profile gets the public profile information about a Spotify user.
package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"

	"golang.org/x/oauth2/clientcredentials"

	"github.com/zmb3/spotify"
)

var userID = flag.String("user", "", "the Spotify user ID to look up")

func main() {
	flag.Parse()

	if *userID == "" {
		fmt.Fprintf(os.Stderr, "Error: missing user ID\n")
		flag.Usage()
		return
	}

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
	user, err := client.GetUsersPublicProfile(spotify.ID(*userID))
	if err != nil {
		fmt.Fprintf(os.Stderr, err.Error())
		return
	}

	fmt.Println("User ID:", user.ID)
	fmt.Println("Display name:", user.DisplayName)
	fmt.Println("Spotify URI:", string(user.URI))
	fmt.Println("Endpoint:", user.Endpoint)
	fmt.Println("Followers:", user.Followers.Count)
}
