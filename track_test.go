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
