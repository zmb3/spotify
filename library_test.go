package spotify

import (
"context"
"encoding/json"
"errors"
"io"
"net/http"
"testing"
)

func TestSaveToLibrary(t *testing.T) {
client, server := testClientString(http.StatusOK, "", func(req *http.Request) {
if req.Method != "PUT" {
t.Errorf("Expected PUT, got %s", req.Method)
}
body, err := io.ReadAll(req.Body)
if err != nil {
t.Fatal(err)
}
var b struct {
URIs []URI `json:"uris"`
}
if err := json.Unmarshal(body, &b); err != nil {
t.Fatal(err)
}
if len(b.URIs) != 2 {
t.Errorf("Expected 2 URIs, got %d", len(b.URIs))
}
})
defer server.Close()

err := client.SaveToLibrary(context.Background(), "spotify:track:4iV5W9uYEdYUVa79Axb7Rh", "spotify:album:1301WleyT98MSxVHPZCA6M")
if err != nil {
t.Error(err)
}
}

func TestSaveToLibraryFailure(t *testing.T) {
client, server := testClientString(http.StatusUnauthorized, `
{
  "error": {
    "status": 401,
    "message": "Invalid access token"
  }
}`)
defer server.Close()
err := client.SaveToLibrary(context.Background(), "spotify:track:4iV5W9uYEdYUVa79Axb7Rh")
if err == nil {
t.Error("Expected error and didn't get one")
}
}

func TestSaveToLibraryWithContextCancelled(t *testing.T) {
client, server := testClientString(http.StatusOK, "")
defer server.Close()

ctx, done := context.WithCancel(context.Background())
done()

err := client.SaveToLibrary(ctx, "spotify:track:4iV5W9uYEdYUVa79Axb7Rh")
if !errors.Is(err, context.Canceled) {
t.Error("Expected error and didn't get one")
}
}

func TestRemoveFromLibrary(t *testing.T) {
client, server := testClientString(http.StatusOK, "", func(req *http.Request) {
if req.Method != "DELETE" {
t.Errorf("Expected DELETE, got %s", req.Method)
}
})
defer server.Close()

err := client.RemoveFromLibrary(context.Background(), "spotify:track:4iV5W9uYEdYUVa79Axb7Rh", "spotify:album:1301WleyT98MSxVHPZCA6M")
if err != nil {
t.Error(err)
}
}

func TestUserHasSavedItems(t *testing.T) {
client, server := testClientString(http.StatusOK, `[ false, true ]`)
defer server.Close()

contains, err := client.UserHasSavedItems(context.Background(), "spotify:track:0udZHhCi7p1YzMlvI4fXoK", "spotify:track:55nlbqqFVnSsArIeYSQlqx")
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

func TestSaveToLibraryEmpty(t *testing.T) {
client, server := testClientString(http.StatusOK, "")
defer server.Close()

err := client.SaveToLibrary(context.Background())
if err == nil {
t.Error("Expected error for empty URIs")
}
}
