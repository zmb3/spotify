package spotify

import (
	"net/http"
	"testing"
)

func TestTransferPlaybackDeviceUnavailable(t *testing.T) {
	client, server := testClientString(http.StatusAccepted, "")
	defer server.Close()
	err := client.TransferPlayback("newdevice", false)
	if err == nil {
		t.Error("expected error since auto retry is disabled")
	}
}

func TestTransferPlayback(t *testing.T) {
	client, server := testClientString(http.StatusNoContent, "")
	defer server.Close()

	err := client.TransferPlayback("newdevice", true)
	if err != nil {
		t.Error(err)
	}
}

func TestVolume(t *testing.T) {
	client, server := testClientString(http.StatusNoContent, "")
	defer server.Close()

	err := client.Volume(50)
	if err != nil {
		t.Error(err)
	}
}

func TestPlayerDevices(t *testing.T) {
	client, server := testClientFile(http.StatusOK, "test_data/player_available_devices.txt")
	defer server.Close()

	list, err := client.PlayerDevices()
	if err != nil {
		t.Error(err)
		return
	}
	if len(list) != 2 {
		t.Error("Expected two devices")
	}

	if list[0].Volume != 100 {
		t.Error("Expected volume to be 100%")
	}
	if list[1].Volume != 0 {
		t.Error("Expected null becomes 0")
	}
}

func TestPlayerState(t *testing.T) {
	client, server := testClientFile(http.StatusOK, "test_data/player_state.txt")
	defer server.Close()

	state, err := client.PlayerState()
	if err != nil {
		t.Error(err)
		return
	}

	if len(state.PlaybackContext.ExternalURLs) != 1 {
		t.Error("Expected one external url")
	}

	if state.Item == nil {
		t.Error("Expected item to be a track")
	}

	if state.Timestamp != 1491302708055 {
		t.Error("Expected timestamp to be 1491302708055")
	}

	if state.Progress != 102509 {
		t.Error("Expected progress to be 102509")
	}

	if state.Playing {
		t.Error("Expected not to be playing")
	}
}

func TestPlayerCurrentlyPlaying(t *testing.T) {
	client, server := testClientFile(http.StatusOK, "test_data/player_currently_playing.txt")
	defer server.Close()

	state, err := client.PlayerCurrentlyPlaying()
	if err != nil {
		t.Error(err)
		return
	}

	if len(state.PlaybackContext.ExternalURLs) != 1 {
		t.Error("Expected one external url")
	}

	if state.Item == nil {
		t.Error("Expected item to be a track")
	}

	if state.Timestamp != 1491302708055 {
		t.Error("Expected timestamp to be 1491302708055")
	}

	if state.Progress != 102509 {
		t.Error("Expected progress to be 102509")
	}

	if state.Playing {
		t.Error("Expected not to be playing")
	}
}

func TestPlayerRecentlyPlayed(t *testing.T) {
	client, server := testClientFile(http.StatusOK, "test_data/player_recently_played.txt")
	defer server.Close()

	items, err := client.PlayerRecentlyPlayed()
	if err != nil {
		t.Fatal(err)
	}

	if len(items) != 20 {
		t.Error("Too few or too many items were returned")
	}

	actualTimePhrase := items[0].PlayedAt.Format("2006-01-02T15:04:05.999Z")
	expectedTimePhrase := "2017-05-27T20:07:54.721Z"

	if actualTimePhrase != expectedTimePhrase {
		t.Errorf("Time of first track was not parsed correctly: [%s] != [%s]", actualTimePhrase, expectedTimePhrase)
	}
}

func TestPlayArgsError(t *testing.T) {
	json := `{
		"error" : {
			"status" : 400,
			"message" : "Only one of either \"context_uri\" or \"uris\" can be specified"
		}
	}`
	client, server := testClientString(http.StatusUnauthorized, json)
	defer server.Close()

	err := client.Play()
	if err == nil {
		t.Error("Expected an error")
	}
}
