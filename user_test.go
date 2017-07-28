package spotify

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

const userResponse = `
{
  "display_name" : "Ronald Pompa",
  "external_urls" : {
    "spotify" : "https://open.spotify.com/user/wizzler"
    },
    "followers" : {
      "href" : null,
      "total" : 3829
    },
    "href" : "https://api.spotify.com/v1/users/wizzler",
    "id" : "wizzler",
    "images" : [ {
      "height" : null,
      "url" : "http://profile-images.scdn.co/images/userprofile/default/9d51820e73667ea5f1e97ea601cf0593b558050e",
      "width" : null
    } ],
    "type" : "user",
    "uri" : "spotify:user:wizzler"
}`

func TestUserProfile(t *testing.T) {
	client, server := testClientString(http.StatusOK, userResponse)
	defer server.Close()

	user, err := client.GetUsersPublicProfile("wizzler")
	if err != nil {
		t.Error(err)
		return
	}
	if user.ID != "wizzler" {
		t.Error("Expected user wizzler, got ", user.ID)
	}
	if f := user.Followers.Count; f != 3829 {
		t.Errorf("Expected 3829 followers, got %d\n", f)
	}
}

func TestCurrentUser(t *testing.T) {
	json := `{
		"country" : "US",
		"display_name" : null,
		"email" : "username@domain.com",
		"external_urls" : {
			"spotify" : "https://open.spotify.com/user/username"
		},
		"followers" : {
			"href" : null,
			"total" : 0
		},
		"href" : "https://api.spotify.com/v1/users/userame",
		"id" : "username",
		"images" : [ ],
		"product" : "premium",
		"type" : "user",
		"uri" : "spotify:user:username",
		"birthdate" : "1985-05-01"
	}`
	client, server := testClientString(http.StatusOK, json)
	defer server.Close()

	me, err := client.CurrentUser()
	if err != nil {
		t.Error(err)
		return
	}
	if me.Country != CountryUSA ||
		me.Email != "username@domain.com" ||
		me.Product != "premium" {
		t.Error("Received incorrect response")
	}
	if me.Birthdate != "1985-05-01" {
		t.Errorf("Expected '1985-05-01', got '%s'\n", me.Birthdate)
	}
}

func TestFollowUsersMissingScope(t *testing.T) {
	json := `{
		"error": {
			"status": 403,
			"message": "Insufficient client scope"
		}
	}`
	client, server := testClientString(http.StatusForbidden, json, func(req *http.Request) {
		if req.URL.Query().Get("type") != "user" {
			t.Error("Request made with the wrong type parameter")
		}
	})
	defer server.Close()

	err := client.FollowUser(ID("exampleuser01"))
	serr, ok := err.(Error)
	if !ok {
		t.Fatal("Expected insufficient client scope error")
	}
	if serr.Status != http.StatusForbidden {
		t.Error("Expected HTTP 403")
	}
}

func TestFollowArtist(t *testing.T) {
	client, server := testClientString(http.StatusNoContent, "", func(req *http.Request) {
		if req.URL.Query().Get("type") != "artist" {
			t.Error("Request made with the wrong type parameter")
		}
	})
	defer server.Close()

	if err := client.FollowArtist("3ge4xOaKvWfhRwgx0Rldov"); err != nil {
		t.Error(err)
	}
}

func TestFollowArtistAutoRetry(t *testing.T) {
	t.Parallel()
	handlers := []http.HandlerFunc{
		// first attempt fails
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Retry-After", "2")
			w.WriteHeader(rateLimitExceededStatusCode)
			io.WriteString(w, `{ "error": { "message": "slow down", "status": 429 } }`)
		}),
		// next attempt succeeds
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNoContent)
		}),
	}

	i := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handlers[i](w, r)
		i++
	}))
	defer server.Close()

	client := &Client{http: http.DefaultClient, baseURL: server.URL + "/", AutoRetry: true}
	if err := client.FollowArtist("3ge4xOaKvWfhRwgx0Rldov"); err != nil {
		t.Error(err)
	}
}

func TestFollowUsersInvalidToken(t *testing.T) {
	json := `{
		"error": {
			"status": 401,
			"message": "Invalid access token"
		}
	}`
	client, server := testClientString(http.StatusUnauthorized, json, func(req *http.Request) {
		if req.URL.Query().Get("type") != "user" {
			t.Error("Request made with the wrong type parameter")
		}
	})
	defer server.Close()

	err := client.FollowUser(ID("dummyID"))
	serr, ok := err.(Error)
	if !ok {
		t.Fatal("Expected invalid token error")
	}
	if serr.Status != http.StatusUnauthorized {
		t.Error("Expected HTTP 401")
	}
}

func TestUserFollows(t *testing.T) {
	json := "[ false, true ]"
	client, server := testClientString(http.StatusOK, json)
	defer server.Close()

	follows, err := client.CurrentUserFollows("artist", ID("74ASZWbe4lXaubB36ztrGX"), ID("08td7MxkoHQkXnWAYD8d6Q"))
	if err != nil {
		t.Error(err)
		return
	}
	if len(follows) != 2 || follows[0] || !follows[1] {
		t.Error("Incorrect result", follows)
	}
}

func TestCurrentUsersTracks(t *testing.T) {
	client, server := testClientFile(http.StatusOK, "test_data/current_users_tracks.txt")
	defer server.Close()

	tracks, err := client.CurrentUsersTracks()
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
	if len(tracks.Tracks) != tracks.Total {
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

	albums, err := client.CurrentUsersAlbums()
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
	if len(albums.Albums) != albums.Total {
		t.Error("Didn't get expected number of results")
		return
	}
	expected := "Love In The Future"
	if albums.Albums[0].Name != expected {
		t.Errorf("Expected '%s', got '%s'\n", expected, albums.Albums[0].Name)
		fmt.Printf("\n%#v\n", albums.Albums[0])
	}

	upc := "886444160742"
	u, ok := albums.Albums[0].ExternalIDs["upc"]
	if !ok {
		t.Error("External IDs missing UPC")
	}
	if u != upc {
		t.Errorf("Wrong UPC: want %s, got %s\n", upc, u)
	}
}

func TestCurrentUsersPlaylists(t *testing.T) {
	client, server := testClientFile(http.StatusOK, "test_data/current_users_playlists.txt")
	defer server.Close()

	playlists, err := client.CurrentUsersPlaylists()
	if err != nil {
		t.Error(err)
	}
	if playlists.Limit != 20 {
		t.Errorf("Expected limit 20, got %d\n", playlists.Limit)
	}
	if playlists.Total != 4 {
		t.Errorf("Expected 4 playlists, got %d\n", playlists.Total)
	}
	tests := []struct {
		Name       string
		Public     bool
		TrackCount uint
	}{
		{"Discover Weekly", false, 30},
		{"Your Favorite Coffeehouse", false, 69},
		{"Afternoon Acoustic", false, 99},
		{"Yoga and Meditation", true, 31},
	}
	for i := range tests {
		p := playlists.Playlists[i]
		if p.Name != tests[i].Name {
			t.Errorf("Expected '%s', got '%s'\n", tests[i].Name, p.Name)
		}
		if p.IsPublic != tests[i].Public {
			t.Errorf("Expected public to be %#v, got %#v\n", tests[i].Public, p.IsPublic)
		}
		if p.Tracks.Total != tests[i].TrackCount {
			t.Errorf("Expected %d tracks, got %d\n", tests[i].TrackCount, p.Tracks.Total)
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

	artists, err := client.CurrentUsersFollowedArtists()
	if err != nil {
		t.Fatal(err)
	}
	exp := 20
	if artists.Limit != exp {
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

	client.CurrentUsersFollowedArtistsOpt(50, "0aV6DOiouImYTqrR5YlIqx")
}
