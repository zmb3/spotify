package spotify

import (
	"fmt"
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

func TestZacFeatured(t *testing.T) {
	t.Skip("skipping real auth test")
	var c Client
	c.AccessToken = "BQAve8rG6wz28MdIsPhuX3v_ziKeaSmBtcE2ncfq3hNn5ypTVTuOD_-Ki7go_qvzgS0Eq_zXO-AbRhmbzmY4t7xBK0mP4wqdTLqLuBJQ3fowQamWjzGL0VJI0I0A2EOpaDQ_wZ33xrVUenJ4eGfZrlWhuM3DI6tI8jXof0tdbiyqyX-oev4aDAZ4AFJS1YEr37Hjdu6qDQ"
	c.TokenType = BearerToken
	msg, p, err := c.FeaturedPlaylists()
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Println(msg)
	for _, item := range p.Playlists {
		fmt.Println(item.Name)
	}
}
