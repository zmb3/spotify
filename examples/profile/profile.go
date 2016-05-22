// Command profile gets the public profile information about a Spotify user.
package main

import (
	"flag"
	"fmt"
	"os"

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

	user, err := spotify.GetUsersPublicProfile(spotify.ID(*userID))
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
