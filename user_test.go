package spotify

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"testing"
)

func TestCurrentUser(t *testing.T) {
	json := `{
		"display_name" : null,
		"external_urls" : {
			"spotify" : "https://open.spotify.com/user/username"
		},
		"href" : "https://api.spotify.com/v1/users/userame",
		"id" : "username",
		"images" : [ ],
		"type" : "user",
		"uri" : "spotify:user:username",
		"birthdate" : "1985-05-01"
	}`
	client, server := testClientString(http.StatusOK, json)
	defer server.Close()

	me, err := client.CurrentUser(context.Background())
	if err != nil {
		t.Error(err)
		return
	}
	if me.Birthdate != "1985-05-01" {
		t.Errorf("Expected '1985-05-01', got '%s'\n", me.Birthdate)
	}
}

func TestCurrentUsersTracks(t *testing.T) {
	client, server := testClientFile(http.StatusOK, "test_data/current_users_tracks.txt")
	defer server.Close()

	tracks, err := client.CurrentUsersTracks(context.Background())
	if err != nil {
		t.Error(err)
		return
	}
	if tracks.Limit != 20 {
		t.Errorf("Expected limit 20, got %d\n", tracks.Limit)
	}
	if tracks.Endpoint != "https://api.spotify.com/v1/me/tracks?offset=0&limit=20" {
		t.Error("Endpoint incorrect")
	}
	if tracks.Total != 3 {
		t.Errorf("Expect 3 results, got %d\n", tracks.Total)
		return
	}
	if len(tracks.Tracks) != int(tracks.Total) {
		t.Error("Didn't get expected number of results")
		return
	}
	expected := "You & I (Nobody In The World)"
	if tracks.Tracks[0].Name != expected {
		t.Errorf("Expected '%s', got '%s'\n", expected, tracks.Tracks[0].Name)
		fmt.Printf("\n%#v\n", tracks.Tracks[0])
	}
}

func TestCurrentUsersAlbums(t *testing.T) {
	client, server := testClientFile(http.StatusOK, "test_data/current_users_albums.txt")
	defer server.Close()

	albums, err := client.CurrentUsersAlbums(context.Background())
	if err != nil {
		t.Error(err)
		return
	}
	if albums.Limit != 20 {
		t.Errorf("Expected limit 20, got %d\n", albums.Limit)
	}
	if albums.Endpoint != "https://api.spotify.com/v1/me/albums?offset=0&limit=20" {
		t.Error("Endpoint incorrect")
	}
	if albums.Total != 2 {
		t.Errorf("Expect 2 results, got %d\n", albums.Total)
		return
	}
	if len(albums.Albums) != int(albums.Total) {
		t.Error("Didn't get expected number of results")
		return
	}
	expected := "Love In The Future"
	if albums.Albums[0].Name != expected {
		t.Errorf("Expected '%s', got '%s'\n", expected, albums.Albums[0].Name)
		fmt.Printf("\n%#v\n", albums.Albums[0])
	}
}

func TestCurrentUsersPlaylists(t *testing.T) {
	client, server := testClientFile(http.StatusOK, "test_data/current_users_playlists.txt")
	defer server.Close()

	playlists, err := client.CurrentUsersPlaylists(context.Background())
	if err != nil {
		t.Error(err)
	}
	if playlists.Limit != 5 {
		t.Errorf("Expected limit 5, got %d\n", playlists.Limit)
	}
	if playlists.Offset != 20 {
		t.Errorf("Expected offset 20, got %d\n", playlists.Offset)
	}
	if playlists.Total != 42 {
		t.Errorf("Expected 42 playlists, got %d\n", playlists.Total)
	}

	tests := []struct {
		Name        string
		Description string
		Public      bool
		TrackCount  int
	}{
		{"Core", "", true, 3},
		{"Black/Atmo/Prog ?", "This is kinda fuzzy", true, 10},
		{"Troll", "", false, 7},
		{"Melomiel", "", true, 3},
		{"HEAVY MIEL", "Deathcore,techdeath, tout ce qui tape plus fort que du melodeath", true, 10},
	}
	for i := range tests {
		p := playlists.Playlists[i]
		if p.Name != tests[i].Name {
			t.Errorf("Expected '%s', got '%s'\n", tests[i].Name, p.Name)
		}
		if p.Description != tests[i].Description {
			t.Errorf("Expected '%s', got '%s'\n", tests[i].Description, p.Description)
		}
		if p.IsPublic != tests[i].Public {
			t.Errorf("Expected public to be %#v, got %#v\n", tests[i].Public, p.IsPublic)
		}
		if int(p.Items.Total) != tests[i].TrackCount {
			t.Errorf("Expected %d tracks, got %d\n", tests[i].TrackCount, p.Items.Total)
		}
	}
}

func TestUsersFollowedArtists(t *testing.T) {
	json := `
{
  "artists" : {
    "items" : [ {
      "external_urls" : {
        "spotify" : "https://open.spotify.com/artist/0I2XqVXqHScXjHhk6AYYRe"
      },
      "followers" : {
        "href" : null,
        "total" : 7753
      },
      "genres" : [ "swedish hip hop" ],
      "href" : "https://api.spotify.com/v1/artists/0I2XqVXqHScXjHhk6AYYRe",
      "id" : "0I2XqVXqHScXjHhk6AYYRe",
      "images" : [ {
        "height" : 640,
        "url" : "https://i.scdn.co/image/2c8c0cea05bf3d3c070b7498d8d0b957c4cdec20",
        "width" : 640
      }, {
        "height" : 300,
        "url" : "https://i.scdn.co/image/394302b42c4b894786943e028cdd46d7baaa29b7",
        "width" : 300
      }, {
        "height" : 64,
        "url" : "https://i.scdn.co/image/ca9df7225ade6e5dfc62e7076709ca3409a7cbbf",
        "width" : 64
      } ],
      "name" : "Afasi & Filthy",
      "popularity" : 54,
      "type" : "artist",
      "uri" : "spotify:artist:0I2XqVXqHScXjHhk6AYYRe"
   } ],
  "next" : "https://api.spotify.com/v1/users/thelinmichael/following?type=artist&after=0aV6DOiouImYTqrR5YlIqx&limit=20",
  "total" : 183,
    "cursors" : {
      "after" : "0aV6DOiouImYTqrR5YlIqx"
    },
   "limit" : 20,
   "href" : "https://api.spotify.com/v1/users/thelinmichael/following?type=artist&limit=20"
  }
}`
	client, server := testClientString(http.StatusOK, json)
	defer server.Close()

	artists, err := client.CurrentUsersFollowedArtists(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	exp := 20
	if int(artists.Limit) != exp {
		t.Errorf("Expected limit %d, got %d\n", exp, artists.Limit)
	}
	if a := artists.Cursor.After; a != "0aV6DOiouImYTqrR5YlIqx" {
		t.Error("Invalid 'after' cursor")
	}
	if l := len(artists.Artists); l != 1 {
		t.Fatalf("Expected 1 artist, got %d\n", l)
	}
	if n := artists.Artists[0].Name; n != "Afasi & Filthy" {
		t.Error("Got wrong artist name")
	}
}

func TestCurrentUsersFollowedArtistsOpt(t *testing.T) {
	client, server := testClientString(http.StatusOK, "{}", func(req *http.Request) {
		if url := req.URL.String(); !strings.HasSuffix(url, "me/following?after=0aV6DOiouImYTqrR5YlIqx&limit=50&type=artist") {
			t.Error("invalid request url")
		}
	})
	defer server.Close()

	_, err := client.CurrentUsersFollowedArtists(context.Background(), Limit(50), After("0aV6DOiouImYTqrR5YlIqx"))
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

func TestCurrentUsersTopArtists(t *testing.T) {
	client, server := testClientFile(http.StatusOK, "test_data/current_users_top_artists.txt")
	defer server.Close()

	artists, err := client.CurrentUsersTopArtists(context.Background())
	if err != nil {
		t.Error(err)
	}
	if artists.Endpoint != "https://api.spotify.com/v1/me/top/artists" {
		t.Error("Endpoint incorrect")
	}
	if artists.Limit != 20 {
		t.Errorf("Expected limit 20, got %d\n", artists.Limit)
	}
	if artists.Total != 10 {
		t.Errorf("Expected total 10, got %d\n", artists.Total)
		return
	}
	if len(artists.Artists) != int(artists.Total) {
		t.Error("Didn't get expected number of results")
		return
	}

	name := "insaneintherainmusic"
	if artists.Artists[0].Name != name {
		t.Errorf("Expected '%s', got '%s'\n", name, artists.Artists[0].Name)
		fmt.Printf("\n%#v\n", artists.Artists[0])
	}
}

func TestCurrentUsersTopTracks(t *testing.T) {
	client, server := testClientFile(http.StatusOK, "test_data/current_users_top_tracks.txt")
	defer server.Close()

	tracks, err := client.CurrentUsersTopTracks(context.Background())
	if err != nil {
		t.Error(err)
	}
	if tracks.Endpoint != "https://api.spotify.com/v1/me/top/tracks" {
		t.Error("Endpoint incorrect")
	}
	if tracks.Limit != 20 {
		t.Errorf("Expected limit 20, got %d\n", tracks.Limit)
	}
	if tracks.Total != 380 {
		t.Errorf("Expected total 380, got %d\n", tracks.Total)
		return
	}
	if len(tracks.Tracks) != int(tracks.Limit) {
		t.Errorf("Didn't get expected number of results")
		return
	}

	name := "Adventure Awaits! (Alola Region Theme)"
	if tracks.Tracks[0].Name != name {
		t.Errorf("Expected '%s', got '%s'\n", name, tracks.Tracks[0].Name)
		fmt.Printf("\n%#v\n", tracks.Tracks[0])
	}
}
