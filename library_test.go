package spotify

import (
	"net/http"
	"testing"
)

func TestUserHasTracks(t *testing.T) {
	client := testClientString(http.StatusOK, `[ false, true ]`)
	addDummyAuth(client)
	contains, err := client.UserHasTracks("0udZHhCi7p1YzMlvI4fXoK", "55nlbqqFVnSsArIeYSQlqx")
	if err != nil {
		t.Error(err)
	}
	if l := len(contains); l != 2 {
		t.Error("Expected 2 results, got", l)
	}
	if contains[0] || !contains[1] {
		t.Error("Expected [false, true], got", contains)
	}
}

func TestAddTracksToLibrary(t *testing.T) {
	client := testClientString(http.StatusOK, "")
	addDummyAuth(client)
	err := client.AddTracksToLibrary("4iV5W9uYEdYUVa79Axb7Rh", "1301WleyT98MSxVHPZCA6M")
	if err != nil {
		t.Error(err)
	}
}

func TestAddTracksToLibraryFailure(t *testing.T) {
	client := testClientString(http.StatusUnauthorized, `
{
  "error": {
    "status": 401,
    "message": "Invalid access token"
  }
}`)
	addDummyAuth(client)
	err := client.AddTracksToLibrary("4iV5W9uYEdYUVa79Axb7Rh", "1301WleyT98MSxVHPZCA6M")
	if err == nil {
		t.Error("Expected error and didn't get one")
	}
}

func TestRemoveTracksFromLibrary(t *testing.T) {
	client := testClientString(http.StatusOK, "")
	addDummyAuth(client)
	err := client.RemoveTracksFromLibrary("4iV5W9uYEdYUVa79Axb7Rh", "1301WleyT98MSxVHPZCA6M")
	if err != nil {
		t.Error(err)
	}
}
