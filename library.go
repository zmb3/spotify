package spotify

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

// SaveToLibrary saves one or more items to the current user's library.
// This call accepts Spotify URIs (e.g. "spotify:track:xxx", "spotify:album:xxx",
// "spotify:artist:xxx", "spotify:playlist:xxx").
//
// Appropriate scopes need to be passed depending on the entities being saved.
func (c *Client) SaveToLibrary(ctx context.Context, uris ...URI) error {
	if l := len(uris); l == 0 {
		return fmt.Errorf("spotify: at least one URI is required")
	}
	spotifyURL := c.baseURL + "me/library"
	body := struct {
		URIs []URI `json:"uris"`
	}{URIs: uris}
	bodyJSON, err := json.Marshal(body)
	if err != nil {
		return err
	}
	req, err := http.NewRequestWithContext(ctx, "PUT", spotifyURL, bytes.NewReader(bodyJSON))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	return c.execute(req, nil)
}

// RemoveFromLibrary removes one or more items from the current user's library.
// This call accepts Spotify URIs (e.g. "spotify:track:xxx", "spotify:album:xxx",
// "spotify:artist:xxx", "spotify:playlist:xxx").
//
// Appropriate scopes need to be passed depending on the entities being removed.
func (c *Client) RemoveFromLibrary(ctx context.Context, uris ...URI) error {
	if l := len(uris); l == 0 {
		return fmt.Errorf("spotify: at least one URI is required")
	}
	spotifyURL := c.baseURL + "me/library"
	body := struct {
		URIs []URI `json:"uris"`
	}{URIs: uris}
	bodyJSON, err := json.Marshal(body)
	if err != nil {
		return err
	}
	req, err := http.NewRequestWithContext(ctx, "DELETE", spotifyURL, bytes.NewReader(bodyJSON))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	return c.execute(req, nil)
}

// UserHasSavedItems checks if one or more items are saved in the current user's library.
// This call accepts Spotify URIs (e.g. "spotify:track:xxx", "spotify:album:xxx").
//
// The result is returned as a slice of bool values in the same order
// in which the URIs were specified.
func (c *Client) UserHasSavedItems(ctx context.Context, uris ...URI) ([]bool, error) {
	if l := len(uris); l == 0 {
		return nil, fmt.Errorf("spotify: at least one URI is required")
	}
	uriStrings := make([]string, len(uris))
	for i, u := range uris {
		uriStrings[i] = string(u)
	}
	spotifyURL := fmt.Sprintf("%sme/library/contains?uris=%s", c.baseURL, strings.Join(uriStrings, ","))

	var result []bool
	err := c.get(ctx, spotifyURL, &result)
	if err != nil {
		return nil, err
	}

	return result, nil
}
