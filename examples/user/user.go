package main

import (
	"fmt"
	"github.com/zmb3/spotify"
	"log"
	"net/http"
)

const redirectURI = "http://localhost:8080/callback"

var (
	auth  = spotify.NewAuthenticator(redirectURI, spotify.ScopeUserLibraryRead)
	ch    = make(chan *spotify.Client)
	state = "abc123"
)

func main() {
	http.HandleFunc("/callback", completeAuth)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Println("Got request for:", r.URL.String())
	})
	go http.ListenAndServe(":8080", nil)

	url := auth.AuthURL(state)
	fmt.Println("Please log in to Spotify by visiting the following page in your browser:", url)

	client := <-ch

	// get all albums saved in library
	var albums []spotify.SavedAlbum
	limit := 50
	offset := 0
	result, _ := client.CurrentUsersAlbumsOpt(&spotify.Options{Limit: &limit, Offset: &offset})
	albums = append(albums, result.Albums...)
	for {
		if err := client.NextPage(result); err != nil {
			break
		}
		albums = append(albums, result.Albums...)
	}

	// print result
	for _, album := range albums {
		fmt.Println(album.Name)
	}
}

func completeAuth(w http.ResponseWriter, r *http.Request) {
	tok, err := auth.Token(state, r)
	if err != nil {
		http.Error(w, "Couldn't get token", http.StatusForbidden)
		log.Fatal(err)
	}
	if st := r.FormValue("state"); st != state {
		http.NotFound(w, r)
		log.Fatalf("State mismatch: %s != %s\n", st, state)
	}
	client := auth.NewClient(tok)

	fmt.Fprintf(w, "Login Completed!")
	ch <- &client
}
