package spotify

import (
	"net/http"
	"os"
	"testing"
)

func TestSearchArtist(t *testing.T) {
	client := testClientFile(http.StatusOK, "test_data/search_artist.txt")
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
	client := testClientFile(http.StatusOK, "test_data/search_tracks.txt")
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

func TestSearchTrackWithFilter(t *testing.T) {
	if os.Getenv("FULLTEST") == "" {
		t.Skip()
		return
	}

	result, err := Search("uptown artist:bruno mars", SearchTypeTrack)
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
	client := testClientFile(http.StatusOK, "test_data/search_trackplaylist.txt")
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
	// we expect to get ErrNoMorePages when trying to get the prev/next page
	// under either of these conditions:

	//  1) there are no results (nil)
	nilResults := &SearchResult{nil, nil, nil, nil}
	if DefaultClient.NextAlbumResults(nilResults) != ErrNoMorePages ||
		DefaultClient.NextArtistResults(nilResults) != ErrNoMorePages ||
		DefaultClient.NextPlaylistResults(nilResults) != ErrNoMorePages ||
		DefaultClient.NextTrackResults(nilResults) != ErrNoMorePages {
		t.Error("Next search result page should have failed for nil results")
	}
	if DefaultClient.PreviousAlbumResults(nilResults) != ErrNoMorePages ||
		DefaultClient.PreviousArtistResults(nilResults) != ErrNoMorePages ||
		DefaultClient.PreviousPlaylistResults(nilResults) != ErrNoMorePages ||
		DefaultClient.PreviousTrackResults(nilResults) != ErrNoMorePages {
		t.Error("Previous search result page should have failed for nil results")
	}
	//  2) the prev/next URL is empty
	emptyURL := &SearchResult{
		Artists:   new(FullArtistPage),
		Albums:    new(SimpleAlbumPage),
		Playlists: new(SimplePlaylistPage),
		Tracks:    new(FullTrackPage),
	}
	if DefaultClient.NextAlbumResults(emptyURL) != ErrNoMorePages ||
		DefaultClient.NextArtistResults(emptyURL) != ErrNoMorePages ||
		DefaultClient.NextPlaylistResults(emptyURL) != ErrNoMorePages ||
		DefaultClient.NextTrackResults(emptyURL) != ErrNoMorePages {
		t.Error("Next search result page should have failed with empty URL")
	}
	if DefaultClient.PreviousAlbumResults(emptyURL) != ErrNoMorePages ||
		DefaultClient.PreviousArtistResults(emptyURL) != ErrNoMorePages ||
		DefaultClient.PreviousPlaylistResults(emptyURL) != ErrNoMorePages ||
		DefaultClient.PreviousTrackResults(emptyURL) != ErrNoMorePages {
		t.Error("Previous search result page should have failed with empty URL")
	}
}

func TestSearchAgainstAPI(t *testing.T) {
	if os.Getenv("FULLTEST") == "" {
		t.Skip()
		return
	}
	t.Parallel()
	res, err := Search("Dave", SearchTypeArtist)
	if err != nil {
		t.Fatal(err)
	}

	// keep requesting the next page of results, up to a maximum of 5 times
	i := 0
	for err = nil; err != ErrNoMorePages && i < 5; err = DefaultClient.NextArtistResults(res) {
		i++
	}
	lastArtist := res.Artists.Artists[0].ID
	// backtrack one page and make sure our artist changed
	if err = DefaultClient.PreviousArtistResults(res); err != nil {
		t.Error(err)
	}
	if lastArtist == res.Artists.Artists[0].ID {
		t.Error("Failed to get previous page")
	}
}
