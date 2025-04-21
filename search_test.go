package spotify

import (
	"context"
	"net/http"
	"testing"
)

func TestSearchArtist(t *testing.T) {
	client, server := testClientFile(http.StatusOK, "test_data/search_artist.txt")
	defer server.Close()

	result, err := client.Search(context.Background(), "tania bowra", SearchTypeArtist)
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

	result, err := client.Search(context.Background(), "uptown", SearchTypeTrack)
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

	result, err := client.Search(context.Background(), "holiday", SearchTypePlaylist|SearchTypeTrack)
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

func TestSearchShow(t *testing.T) {
	client, server := testClientFile(http.StatusOK, "test_data/search_show.txt")
	defer server.Close()

	options := []RequestOption{
		Market("CO"),
	}

	result, err := client.Search(context.Background(), "go time", SearchTypeShow, options...)
	if err != nil {
		t.Error(err)
	}
	if result.Albums != nil {
		t.Error("Searched for shows but received album results")
	}
	if result.Playlists != nil {
		t.Error("Searched for shows but received playlist results")
	}
	if result.Tracks != nil {
		t.Error("Searched for shows but received track results")
	}
	if result.Shows == nil || len(result.Shows.Shows) == 0 {
		t.Error("Didn't receive show results")
	}
	if result.Shows.Shows[0].Name != "Go Time: Golang, Software Engineering" {
		t.Error("Got wrong show name")
	}
}

func TestPrevNextSearchPageErrors(t *testing.T) {
	client, server := testClientString(0, "")
	defer server.Close()

	// we expect to get ErrNoMorePages when trying to get the prev/next page
	// under either of these conditions:

	//  1) there are no results (nil)
	nilResults := &SearchResult{nil, nil, nil, nil, nil, nil}
	if client.NextAlbumResults(context.Background(), nilResults) != ErrNoMorePages ||
		client.NextArtistResults(context.Background(), nilResults) != ErrNoMorePages ||
		client.NextPlaylistResults(context.Background(), nilResults) != ErrNoMorePages ||
		client.NextTrackResults(context.Background(), nilResults) != ErrNoMorePages ||
		client.NextShowResults(context.Background(), nilResults) != ErrNoMorePages ||
		client.NextEpisodeResults(context.Background(), nilResults) != ErrNoMorePages {
		t.Error("Next search result page should have failed for nil results")
	}
	if client.PreviousAlbumResults(context.Background(), nilResults) != ErrNoMorePages ||
		client.PreviousArtistResults(context.Background(), nilResults) != ErrNoMorePages ||
		client.PreviousPlaylistResults(context.Background(), nilResults) != ErrNoMorePages ||
		client.PreviousTrackResults(context.Background(), nilResults) != ErrNoMorePages ||
		client.PreviousShowResults(context.Background(), nilResults) != ErrNoMorePages ||
		client.PreviousEpisodeResults(context.Background(), nilResults) != ErrNoMorePages {
		t.Error("Previous search result page should have failed for nil results")
	}
	//  2) the prev/next URL is empty
	emptyURL := &SearchResult{
		Artists:   new(FullArtistPage),
		Albums:    new(SimpleAlbumPage),
		Playlists: new(SimplePlaylistPage),
		Tracks:    new(FullTrackPage),
		Shows:     new(SimpleShowPage),
		Episodes:  new(SimpleEpisodePage),
	}
	if client.NextAlbumResults(context.Background(), emptyURL) != ErrNoMorePages ||
		client.NextArtistResults(context.Background(), emptyURL) != ErrNoMorePages ||
		client.NextPlaylistResults(context.Background(), emptyURL) != ErrNoMorePages ||
		client.NextTrackResults(context.Background(), emptyURL) != ErrNoMorePages ||
		client.NextShowResults(context.Background(), emptyURL) != ErrNoMorePages ||
		client.NextEpisodeResults(context.Background(), emptyURL) != ErrNoMorePages {
		t.Error("Next search result page should have failed with empty URL")
	}
	if client.PreviousAlbumResults(context.Background(), emptyURL) != ErrNoMorePages ||
		client.PreviousArtistResults(context.Background(), emptyURL) != ErrNoMorePages ||
		client.PreviousPlaylistResults(context.Background(), emptyURL) != ErrNoMorePages ||
		client.PreviousTrackResults(context.Background(), emptyURL) != ErrNoMorePages ||
		client.PreviousShowResults(context.Background(), emptyURL) != ErrNoMorePages ||
		client.PreviousEpisodeResults(context.Background(), emptyURL) != ErrNoMorePages {
		t.Error("Previous search result page should have failed with empty URL")
	}
}
