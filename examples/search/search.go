package main

import (
	"fmt"

	"github.com/zmb3/spotify"
)

func main() {
	// search for "holiday"
	results, err := spotify.Search("holiday", spotify.SearchTypePlaylist|spotify.SearchTypeAlbum)
	if err != nil {
		fmt.Printf("Error occurred: %v\n", err)
		return
	}

	// handle album results
	if results.Albums != nil {
		for _, item := range results.Albums.Albums {
			fmt.Println("Album: ", item.Name)
		}
	}
	// handle playlist results
	if results.Playlists != nil {
		for _, item := range results.Playlists.Playlists {
			fmt.Println("Playlist: ", item.Name)
		}
	}
}
