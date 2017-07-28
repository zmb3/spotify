package spotify

import (
	"net/http"
	"testing"
)

func TestSearchArtist(t *testing.T) {
	client, server := testClientFile(http.StatusOK, "test_data/search_artist.txt")
	defer server.Close()

	result, err := client.Search("tania bowra", SearchTypeArtist)
	if err != nil {
		t.Error(err)
	}
	if result.Albums != nil {
		t.Error("Searched for artists but received album results")
	}
	if result.Playlists != nil {
		t.Error("Searched for artists but received playlist results")
	}
	if result.Tracks != nil {
		t.Error("Searched for artists but received track results")
	}
	if result.Artists == nil || len(result.Artists.Artists) == 0 {
		t.Error("Didn't receive artist results")
	}
	if result.Artists.Artists[0].Name != "Tania Bowra" {
		t.Error("Got wrong artist name")
	}
}

func TestSearchTracks(t *testing.T) {
	client, server := testClientFile(http.StatusOK, "test_data/search_tracks.txt")
	defer server.Close()

	result, err := client.Search("uptown", SearchTypeTrack)
	if err != nil {
		t.Error(err)
	}
	if result.Albums != nil {
		t.Error("Searched for tracks but got album results")
	}
	if result.Playlists != nil {
		t.Error("Searched for tracks but got playlist results")
	}
	if result.Artists != nil {
		t.Error("Searched for tracks but got artist results")
	}
	if result.Tracks == nil || len(result.Tracks.Tracks) == 0 {
		t.Fatal("Didn't receive track results")
	}
	if name := result.Tracks.Tracks[0].Name; name != "Uptown Funk" {
		t.Errorf("Got %s, wanted Uptown Funk\n", name)
	}
}

func TestSearchPlaylistTrack(t *testing.T) {
	client, server := testClientFile(http.StatusOK, "test_data/search_trackplaylist.txt")
	defer server.Close()

	result, err := client.Search("holiday", SearchTypePlaylist|SearchTypeTrack)
	if err != nil {
		t.Error(err)
	}
	if result.Albums != nil {
		t.Error("Searched for playlists and tracks but received album results")
	}
	if result.Artists != nil {
		t.Error("Searched for playlists and tracks but received artist results")
	}
	if result.Tracks == nil {
		t.Error("Didn't receive track results")
	}
	if result.Playlists == nil {
		t.Error("Didn't receive playlist results")
	}
}

func TestPrevNextSearchPageErrors(t *testing.T) {
	client, server := testClientString(0, "")
	defer server.Close()

	// we expect to get ErrNoMorePages when trying to get the prev/next page
	// under either of these conditions:

	//  1) there are no results (nil)
	nilResults := &SearchResult{nil, nil, nil, nil}
	if client.NextAlbumResults(nilResults) != ErrNoMorePages ||
		client.NextArtistResults(nilResults) != ErrNoMorePages ||
		client.NextPlaylistResults(nilResults) != ErrNoMorePages ||
		client.NextTrackResults(nilResults) != ErrNoMorePages {
		t.Error("Next search result page should have failed for nil results")
	}
	if client.PreviousAlbumResults(nilResults) != ErrNoMorePages ||
		client.PreviousArtistResults(nilResults) != ErrNoMorePages ||
		client.PreviousPlaylistResults(nilResults) != ErrNoMorePages ||
		client.PreviousTrackResults(nilResults) != ErrNoMorePages {
		t.Error("Previous search result page should have failed for nil results")
	}
	//  2) the prev/next URL is empty
	emptyURL := &SearchResult{
		Artists:   new(FullArtistPage),
		Albums:    new(SimpleAlbumPage),
		Playlists: new(SimplePlaylistPage),
		Tracks:    new(FullTrackPage),
	}
	if client.NextAlbumResults(emptyURL) != ErrNoMorePages ||
		client.NextArtistResults(emptyURL) != ErrNoMorePages ||
		client.NextPlaylistResults(emptyURL) != ErrNoMorePages ||
		client.NextTrackResults(emptyURL) != ErrNoMorePages {
		t.Error("Next search result page should have failed with empty URL")
	}
	if client.PreviousAlbumResults(emptyURL) != ErrNoMorePages ||
		client.PreviousArtistResults(emptyURL) != ErrNoMorePages ||
		client.PreviousPlaylistResults(emptyURL) != ErrNoMorePages ||
		client.PreviousTrackResults(emptyURL) != ErrNoMorePages {
		t.Error("Previous search result page should have failed with empty URL")
	}
}
