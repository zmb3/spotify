package spotify

import (
	"net/http"
	"testing"
)

func TestFindTrack(t *testing.T) {
	server, client := testClientFromFile(http.StatusOK, "test_data/find_track.txt", t)
	defer server.Close()
	if t.Failed() {
		return
	}
	track, err := client.FindTrack(ID("1zHlj4dQ8ZAtrayhuDDmkY"))
	if err != nil {
		t.Error(err)
		return
	}
	if track.Name != "Timber" {
		t.Errorf("Wanted track Timer, got %s\n", track.Name)
	}
}

func TestFindTracksSimple(t *testing.T) {
	server, client := testClientFromFile(http.StatusOK, "test_data/find_tracks_simple.txt", t)
	defer server.Close()
	if t.Failed() {
		return
	}
	tracks, err := client.FindTracks(ID("0eGsygTp906u18L0Oimnem"), ID("1lDWb6b6ieDQ2xT7ewTC3G"))
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
	server, client := testClientFromFile(http.StatusOK, "test_data/find_tracks_notfound.txt", t)
	defer server.Close()
	if t.Failed() {
		return
	}
	tracks, err := client.FindTracks(ID("0eGsygTp906u18L0Oimnem"), ID("1lDWb6b6iecccdsdckTC3G"))
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
