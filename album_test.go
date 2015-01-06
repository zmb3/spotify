package spotify

import "testing"

// The example from https://developer.spotify.com/web-api/get-album/
func TestFindAlbum(t *testing.T) {
	var c Client
	album, err := c.FindAlbum(ID("0sNOF9WDwhWunNAHPD3Baj"))
	if err != nil || album == nil {
		t.Error(err.Error())
	}
	if album.Name != "She's So Unusual" {
		t.Error("Got wrong album")
	}
}

// The example from https://developer.spotify.com/web-api/get-several-albums/
func TestFindAlbums(t *testing.T) {
	var c Client
	res, err := c.FindAlbums(ID("41MnTivkwTO3UUJ8DrqEJJ"), ID("6JWc4iAiJ9FjyK0B59ABb4"), ID("6UXCm6bOO4gFlDQZV5yL37"))
	if err != nil {
		t.Error(err.Error())
		return
	}
	if len(res.Albums) != 3 {
		t.Error("Expected 3 albums, got " + string(len(res.Albums)))
		return
	}
	expectedAlbums := []string{
		"The Best Of Keane (Deluxe Edition)",
		"Strangeland",
		"Night Train",
	}
	for i, name := range expectedAlbums {
		if res.Albums[i].Name != name {
			t.Error("Expected album", name, "but got", res.Albums[i].Name)
		}
	}
}
