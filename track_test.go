package spotify

import (
	"net/http"
	"testing"
)

func TestFindTrack(t *testing.T) {
	client := testClientFile(http.StatusOK, "test_data/find_track.txt")
	track, err := client.GetTrack(ID("1zHlj4dQ8ZAtrayhuDDmkY"))
	if err != nil {
		t.Error(err)
		return
	}
	if track.Name != "Timber" {
		t.Errorf("Wanted track Timer, got %s\n", track.Name)
	}
}

func TestFindTracksSimple(t *testing.T) {
	client := testClientFile(http.StatusOK, "test_data/find_tracks_simple.txt")
	tracks, err := client.GetTracks(ID("0eGsygTp906u18L0Oimnem"), ID("1lDWb6b6ieDQ2xT7ewTC3G"))
	if err != nil {
		t.Error(err)
		return
	}
	if l := len(tracks); l != 2 {
		t.Errorf("Wanted 2 tracks, got %d\n", l)
		return
	}

}

func TestFindTracksNotFound(t *testing.T) {
	client := testClientFile(http.StatusOK, "test_data/find_tracks_notfound.txt")
	tracks, err := client.GetTracks(ID("0eGsygTp906u18L0Oimnem"), ID("1lDWb6b6iecccdsdckTC3G"))
	if err != nil {
		t.Error(err)
		return
	}
	if l := len(tracks); l != 2 {
		t.Errorf("Expected 2 results, got %d\n", l)
		return
	}
	if tracks[0].Name != "Mr. Brightside" {
		t.Errorf("Expected Mr. Brightside, got %s\n", tracks[0].Name)
	}
	if tracks[1] != nil {
		t.Error("Expected nil track (invalid ID) but got valid track")
	}
}
