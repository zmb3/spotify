package spotify

import (
	"net/http"
	"testing"
)

func TestFeaturedPlaylistNoAuth(t *testing.T) {
	var client Client
	_, _, err := client.FeaturedPlaylists()
	if err == nil {
		t.Error("Call should have failed without authorization:", err)
	}
}

func TestFeaturedPlaylists(t *testing.T) {
	client := testClientFile(http.StatusOK, "test_data/featured_playlists.txt")
	addDummyAuth(client)

	country := "SE"
	opt := PlaylistOptions{}
	opt.Country = &country

	msg, p, err := client.FeaturedPlaylistsOpt(&opt)
	if err != nil {
		t.Error(err)
		return
	}
	if msg != "Enjoy a mellow afternoon." {
		t.Errorf("Want 'Enjoy a mellow afternoon.', got'%s'\n", msg)
		return
	}
	if p.Playlists == nil || len(p.Playlists) == 0 {
		t.Error("Empty playlists result")
		return
	}
	expected := "Hangover Friendly Singer-Songwriter"
	if name := p.Playlists[0].Name; name != expected {
		t.Errorf("Want '%s', got '%s'\n", name, expected)
	}
}

func TestFeaturedPlaylistsExpiredToken(t *testing.T) {
	json := `{
		"error": {
			"status": 401,
			"message": "The access token expired"
		}
	}`
	client := testClientString(http.StatusUnauthorized, json)
	addDummyAuth(client)

	msg, pl, err := client.FeaturedPlaylists()
	if msg != "" || pl != nil || err == nil {
		t.Error("Expected an error")
		return
	}
	serr, ok := err.(Error)
	if !ok {
		t.Error("Expected spotify Error")
		return
	}
	if serr.Status != http.StatusUnauthorized {
		t.Error("Expected HTTP 401")
	}
}
