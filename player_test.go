package spotify

import (
	"context"
	"net/http"
	"testing"
)

func TestTransferPlaybackDeviceUnavailable(t *testing.T) {
	client, server := testClientString(http.StatusNotFound, "")
	defer server.Close()
	err := client.TransferPlayback(context.Background(), "newdevice", false)
	if err == nil {
		t.Error("expected error since auto retry is disabled")
	}
}

func TestTransferPlayback(t *testing.T) {
	client, server := testClientString(http.StatusNoContent, "")
	defer server.Close()

	err := client.TransferPlayback(context.Background(), "newdevice", true)
	if err != nil {
		t.Error(err)
	}
}

func TestVolume(t *testing.T) {
	client, server := testClientString(http.StatusNoContent, "")
	defer server.Close()

	err := client.Volume(context.Background(), 50)
	if err != nil {
		t.Error(err)
	}
}

func TestQueue(t *testing.T) {
	client, server := testClientString(http.StatusNoContent, "")
	defer server.Close()

	err := client.QueueSong(context.Background(), "4JpKVNYnVcJ8tuMKjAj50A")
	if err != nil {
		t.Error(err)
	}
}

func TestPlayerDevices(t *testing.T) {
	client, server := testClientFile(http.StatusOK, "test_data/player_available_devices.txt")
	defer server.Close()

	list, err := client.PlayerDevices(context.Background())
	if err != nil {
		t.Error(err)
		return
	}
	if len(list) != 2 {
		t.Error("Expected two devices")
	}

	if list[0].Volume != 100 {
		t.Error("Expected volume to be 100 percent")
	}
	if list[1].Volume != 0 {
		t.Error("Expected null becomes 0")
	}
}

func TestPlayerState(t *testing.T) {
	client, server := testClientFile(http.StatusOK, "test_data/player_state.txt")
	defer server.Close()

	state, err := client.PlayerState(context.Background())
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

	state, err := client.PlayerCurrentlyPlaying(context.Background())
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

func TestPlayerCurrentlyPlayingEpisode(t *testing.T) {
	client, server := testClientFile(http.StatusOK, "test_data/player_currently_playing_episode.json")
	defer server.Close()

	current, err := client.PlayerCurrentlyPlaying(context.Background())
	if err != nil {
		t.Error(err)
		return
	}

	if current.Item == nil {
		t.Error("Expected item to be a episode")
	}

	expectedName := "300 multiple choices"
	actualName := current.Item.Episode.Name
	if expectedName != actualName {
		t.Errorf("Got '%s', expected '%s'\n", actualName, expectedName)
	}

	if current.Playing {
		t.Error("Expected not to be playing")
	}
}

func TestPlayerCurrentlyPlayingOverride(t *testing.T) {
	var types string
	client, server := testClientString(http.StatusForbidden, "", func(r *http.Request) {
		types = r.URL.Query().Get("additional_types")
	})
	defer server.Close()

	_, _ = client.PlayerCurrentlyPlaying(context.Background(), AdditionalTypes(EpisodeAdditionalType))

	if types != "episode" {
		t.Errorf("Expected additional type episode, got %s\n", types)
	}
}

func TestPlayerCurrentlyPlayingDefault(t *testing.T) {
	var types string
	client, server := testClientString(http.StatusForbidden, "", func(r *http.Request) {
		types = r.URL.Query().Get("additional_types")
	})
	defer server.Close()

	_, _ = client.PlayerCurrentlyPlaying(context.Background())

	if types != "episode,track" {
		t.Errorf("Expected additional type episode,track, got %s\n", types)
	}
}

func TestPlayerRecentlyPlayed(t *testing.T) {
	client, server := testClientFile(http.StatusOK, "test_data/player_recently_played.txt")
	defer server.Close()

	items, err := client.PlayerRecentlyPlayed(context.Background())
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

	actualAlbumName := items[0].Track.Album.Name
	expectedAlbumName := "Immortalized"

	if actualAlbumName != expectedAlbumName {
		t.Errorf("Album name of first track was not parsed correctly: [%s] != [%s]", actualAlbumName, expectedAlbumName)
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

	err := client.Play(context.Background())
	if err == nil {
		t.Error("Expected an error")
	}
}

func TestGetQueue(t *testing.T) {
	client, server := testClientFile(http.StatusOK, "test_data/get_queue.txt")
	defer server.Close()

	queue, err := client.GetQueue(context.Background())
	if err != nil {
		t.Error(err)
	}
	if l := len(queue.Items); l == 0 {
		t.Fatal("Didn't get any results")
	} else if l != 20 {
		t.Errorf("Got %d playlists, expected 20\n", l)
	}

	p := queue.Items[0].SimpleTrack
	if p.Name != "This Is the End (For You My Friend)" {
		t.Error("Expected 'This Is the End (For You My Friend)', got", p.Name)
	}

	p = queue.CurrentlyPlaying.SimpleTrack

	if p.Name != "Know Your Enemy" {
		t.Error("Expected 'Know Your Enemy', got", p.Name)
	}
}
