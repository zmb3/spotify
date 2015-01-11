package spotify

import "testing"

func TestFindArtist(t *testing.T) {
  server, client := testClientFromFile(200, "test_data/find_artist.txt", t)
  defer server.Close()
  if t.Failed() {
    return
  }
  artist, err := client.FindArtist(ID("0TnOYISbd1XYRBk9myaseg"))
  if err != nil {
    t.Error(err)
    return
  }
  if followers := artist.Followers.Count; followers != 2265279 {
    t.Errorf("Got %d followers, want 2265279\n", followers)
    return
  }
  if artist.Name != "Pitbull" {
    t.Error("Got ", artist.Name, ", wanted Pitbull")
  }
}

func TestArtistTopTracks(t *testing.T) {
	server, client := testClientFromFile(200, "test_data/artist_top_tracks.txt", t)
	defer server.Close()
	if t.Failed() {
		return
	}
	tracks, err := client.ArtistsTopTracks(ID("43ZHCT0cAZBISjO8DG9PnE"), "SE")
	if err != nil {
		t.Error(err)
		return
	}
	l := len(tracks)
	if l != 10 {
		t.Errorf("Got %d tracks, expected 10\n", l)
	}
	track := tracks[9]
	if track.Name != "(You're The) Devil in Disguise" {
		t.Error("Incorrect track name")
	}
	if track.TrackNumber != 24 {
		t.Errorf("Track number was %d, expected 24\n", track.TrackNumber)
	}
}
