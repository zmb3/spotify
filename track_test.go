// Copyright 2014, 2015 Zac Bergquist
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package spotify

import (
	"net/http"
	"testing"
)

func TestFindTrack(t *testing.T) {
	client := testClientFile(http.StatusOK, "test_data/find_track.txt")
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
	client := testClientFile(http.StatusOK, "test_data/find_tracks_simple.txt")
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
	client := testClientFile(http.StatusOK, "test_data/find_tracks_notfound.txt")
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
