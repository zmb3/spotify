package spotify

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

// UserHasTracks checks if one or more tracks are saved to the current user's
// "Your Music" library.
//
// Deprecated: Use [Client.UserHasLibraryItems].
func (c *Client) UserHasTracks(ctx context.Context, ids ...ID) ([]bool, error) {
	return c.libraryContains(ctx, "tracks", ids...)
}

// UserHasAlbums checks if one or more albums are saved to the current user's
// "Your Albums" library.
//
// Deprecated: Use [Client.UserHasLibraryItems].
func (c *Client) UserHasAlbums(ctx context.Context, ids ...ID) ([]bool, error) {
	return c.libraryContains(ctx, "albums", ids...)
}

func (c *Client) libraryContains(ctx context.Context, typ string, ids ...ID) ([]bool, error) {
	if l := len(ids); l == 0 || l > 50 {
		return nil, errors.New("spotify: supports 1 to 50 IDs per call")
	}
	spotifyURL := fmt.Sprintf("%sme/%s/contains?ids=%s", c.baseURL, typ, strings.Join(toStringSlice(ids), ","))

	var result []bool

	err := c.get(ctx, spotifyURL, &result)
	if err != nil {
		return nil, err
	}

	return result, err
}

// UserHasLibraryItems checks if one or more items are saved to the current user's
// library.
func (c *Client) UserHasLibraryItems(ctx context.Context, uris ...URI) ([]bool, error) {
	if l := len(uris); l == 0 || l > 40 {
		return nil, errors.New("spotify: supports 1 to 40 URIs per call")
	}
	v := url.Values{}
	v.Set("uris", strings.Join(toURIStringSlice(uris), ","))
	spotifyURL := c.baseURL + "me/library/contains?" + v.Encode()

	var result []bool

	err := c.get(ctx, spotifyURL, &result)
	if err != nil {
		return nil, err
	}

	return result, err
}

// AddTracksToLibrary saves one or more tracks to the current user's
// "Your Music" library.  This call requires the [ScopeUserLibraryModify] scope.
// A track can only be saved once; duplicate IDs are ignored.
//
// Deprecated: Use [Client.AddItemsToLibrary].
func (c *Client) AddTracksToLibrary(ctx context.Context, ids ...ID) error {
	return c.modifyLibrary(ctx, "tracks", true, ids...)
}

// RemoveTracksFromLibrary removes one or more tracks from the current user's
// "Your Music" library.  This call requires the [ScopeUserModifyLibrary] scope.
// Trying to remove a track when you do not have the user's authorization
// results in an [Error] with the status code set to [net/http.StatusUnauthorized].
//
// Deprecated: Use [Client.RemoveItemsFromLibrary].
func (c *Client) RemoveTracksFromLibrary(ctx context.Context, ids ...ID) error {
	return c.modifyLibrary(ctx, "tracks", false, ids...)
}

// AddAlbumsToLibrary saves one or more albums to the current user's
// "Your Albums" library.  This call requires the [ScopeUserLibraryModify] scope.
// A track can only be saved once; duplicate IDs are ignored.
//
// Deprecated: Use [Client.AddItemsToLibrary].
func (c *Client) AddAlbumsToLibrary(ctx context.Context, ids ...ID) error {
	return c.modifyLibrary(ctx, "albums", true, ids...)
}

// RemoveAlbumsFromLibrary removes one or more albums from the current user's
// "Your Albums" library.  This call requires the [ScopeUserModifyLibrary] scope.
// Trying to remove a track when you do not have the user's authorization
// results in an [Error] with the status code set to [net/http.StatusUnauthorized].
//
// Deprecated: Use [Client.RemoveItemsFromLibrary].
func (c *Client) RemoveAlbumsFromLibrary(ctx context.Context, ids ...ID) error {
	return c.modifyLibrary(ctx, "albums", false, ids...)
}

func (c *Client) modifyLibrary(ctx context.Context, typ string, add bool, ids ...ID) error {
	if l := len(ids); l == 0 || l > 50 {
		return errors.New("spotify: this call supports 1 to 50 IDs per call")
	}
	spotifyURL := fmt.Sprintf("%sme/%s?ids=%s", c.baseURL, typ, strings.Join(toStringSlice(ids), ","))
	method := "DELETE"
	if add {
		method = "PUT"
	}
	req, err := http.NewRequestWithContext(ctx, method, spotifyURL, nil)
	if err != nil {
		return err
	}
	return c.execute(req, nil)
}

// AddItemsToLibrary saves one or more items to the current user's library.
func (c *Client) AddItemsToLibrary(ctx context.Context, uris ...URI) error {
	return c.modifyLibraryItems(ctx, true, uris...)
}

// RemoveItemsFromLibrary removes one or more items from the current user's library.
func (c *Client) RemoveItemsFromLibrary(ctx context.Context, uris ...URI) error {
	return c.modifyLibraryItems(ctx, false, uris...)
}

func (c *Client) modifyLibraryItems(ctx context.Context, add bool, uris ...URI) error {
	if l := len(uris); l == 0 || l > 40 {
		return errors.New("spotify: this call supports 1 to 40 URIs per call")
	}
	v := url.Values{}
	v.Set("uris", strings.Join(toURIStringSlice(uris), ","))
	spotifyURL := c.baseURL + "me/library?" + v.Encode()
	method := "DELETE"
	if add {
		method = "PUT"
	}
	req, err := http.NewRequestWithContext(ctx, method, spotifyURL, nil)
	if err != nil {
		return err
	}
	return c.execute(req, nil)
}

func toURIStringSlice(uris []URI) []string {
	result := make([]string, len(uris))
	for i, uri := range uris {
		result[i] = string(uri)
	}
	return result
}
