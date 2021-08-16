package spotify

import (
	"context"
	"net/http"
	"testing"
)

func TestUserHasTracks(t *testing.T) {
	client, server := testClientString(http.StatusOK, `[ false, true ]`)
	defer server.Close()

	contains, err := client.UserHasTracks(context.Background(), "0udZHhCi7p1YzMlvI4fXoK", "55nlbqqFVnSsArIeYSQlqx")
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
	client, server := testClientString(http.StatusOK, "")
	defer server.Close()

	err := client.AddTracksToLibrary(context.Background(), "4iV5W9uYEdYUVa79Axb7Rh", "1301WleyT98MSxVHPZCA6M")
	if err != nil {
		t.Error(err)
	}
}

func TestAddTracksToLibraryFailure(t *testing.T) {
	client, server := testClientString(http.StatusUnauthorized, `
{
  "error": {
    "status": 401,
    "message": "Invalid access token"
  }
}`)
	defer server.Close()
	err := client.AddTracksToLibrary(context.Background(), "4iV5W9uYEdYUVa79Axb7Rh", "1301WleyT98MSxVHPZCA6M")
	if err == nil {
		t.Error("Expected error and didn't get one")
	}
}

func TestRemoveTracksFromLibrary(t *testing.T) {
	client, server := testClientString(http.StatusOK, "")
	defer server.Close()

	err := client.RemoveTracksFromLibrary(context.Background(), "4iV5W9uYEdYUVa79Axb7Rh", "1301WleyT98MSxVHPZCA6M")
	if err != nil {
		t.Error(err)
	}
}

func TestUserHasAlbums(t *testing.T) {
	client, server := testClientString(http.StatusOK, `[ false, true ]`)
	defer server.Close()

	contains, err := client.UserHasAlbums(context.Background(), "0udZHhCi7p1YzMlvI4fXoK", "55nlbqqFVnSsArIeYSQlqx")
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

func TestAddAlbumsToLibrary(t *testing.T) {
	client, server := testClientString(http.StatusOK, "")
	defer server.Close()

	err := client.AddAlbumsToLibrary(context.Background(), "4iV5W9uYEdYUVa79Axb7Rh", "1301WleyT98MSxVHPZCA6M")
	if err != nil {
		t.Error(err)
	}
}

func TestAddAlbumsToLibraryFailure(t *testing.T) {
	client, server := testClientString(http.StatusUnauthorized, `
{
  "error": {
    "status": 401,
    "message": "Invalid access token"
  }
}`)
	defer server.Close()
	err := client.AddAlbumsToLibrary(context.Background(), "4iV5W9uYEdYUVa79Axb7Rh", "1301WleyT98MSxVHPZCA6M")
	if err == nil {
		t.Error("Expected error and didn't get one")
	}
}

func TestRemoveAlbumsFromLibrary(t *testing.T) {
	client, server := testClientString(http.StatusOK, "")
	defer server.Close()

	err := client.RemoveAlbumsFromLibrary(context.Background(), "4iV5W9uYEdYUVa79Axb7Rh", "1301WleyT98MSxVHPZCA6M")
	if err != nil {
		t.Error(err)
	}
}
