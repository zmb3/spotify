package spotify

import (
	"context"
	"net/http"
	"testing"
)

func TestFindTrack(t *testing.T) {
	client, server := testClientFile(http.StatusOK, "test_data/find_track.txt")
	defer server.Close()

	track, err := client.GetTrack(context.Background(), "1zHlj4dQ8ZAtrayhuDDmkY")
	if err != nil {
		t.Error(err)
		return
	}
	if track.Name != "Timber" {
		t.Errorf("Wanted track Timer, got %s\n", track.Name)
	}
}

func TestFindTrackWithFloats(t *testing.T) {
	client, server := testClientFile(http.StatusOK, "test_data/find_track_with_floats.txt")
	defer server.Close()

	track, err := client.GetTrack(context.Background(), "1zHlj4dQ8ZAtrayhuDDmkY")
	if err != nil {
		t.Error(err)
		return
	}
	if track.Name != "Timber" {
		t.Errorf("Wanted track Timer, got %s\n", track.Name)
	}
}
