package spotify

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
)

// UserHasTracks checks if one or more tracks are saved to the current user's
// "Your Music" library.
func (c *Client) UserHasTracks(ids ...ID) ([]bool, error) {
	if l := len(ids); l == 0 || l > 50 {
		return nil, errors.New("spotify: UserHasTracks supports 1 to 50 IDs per call")
	}
	spotifyURL := fmt.Sprintf("%sme/tracks/contains?ids=%s", c.baseURL, strings.Join(toStringSlice(ids), ","))

	var result []bool

	err := c.get(spotifyURL, &result)
	if err != nil {
		return nil, err
	}

	return result, err
}

// AddTracksToLibrary saves one or more tracks to the current user's
// "Your Music" library.  This call requires the ScopeUserLibraryModify scope.
// A track can only be saved once; duplicate IDs are ignored.
func (c *Client) AddTracksToLibrary(ids ...ID) error {
	return c.modifyLibraryTracks(true, ids...)
}

// RemoveTracksFromLibrary removes one or more tracks from the current user's
// "Your Music" library.  This call requires the ScopeUserModifyLibrary scope.
// Trying to remove a track when you do not have the user's authorization
// results in a `spotify.Error` with the status code set to http.StatusUnauthorized.
func (c *Client) RemoveTracksFromLibrary(ids ...ID) error {
	return c.modifyLibraryTracks(false, ids...)
}

func (c *Client) modifyLibraryTracks(add bool, ids ...ID) error {
	if l := len(ids); l == 0 || l > 50 {
		return errors.New("spotify: this call supports 1 to 50 IDs per call")
	}
	spotifyURL := fmt.Sprintf("%sme/tracks?ids=%s", c.baseURL, strings.Join(toStringSlice(ids), ","))
	method := "DELETE"
	if add {
		method = "PUT"
	}
	req, err := http.NewRequest(method, spotifyURL, nil)
	if err != nil {
		return err
	}
	err = c.execute(req, nil)
	if err != nil {
		return err
	}
	return nil
}

// UserSavedTracks returns the tracks which the current user has saved. This 
// call requires the ScopeUserLibraryRead scope. Limit specifies the max items
// per page.
func (c *Client) UserSavedTracks(limit int) (SavedTrackPage, error) {
  return c.UserSavedTracksOpt(limit, 0, "")
}

// UserSavedTracksOpt returns tracks which the current user has saved. This
// call requires the ScopeUserLibraryRead scope. Limit and offset are used for
// manual paging. Limit is the max number of objects to return in the page,
// offset is the index of the first object in the page. Market is an ISO 3166-1
// alpha-2 country code, or the string "from_token". This is used for track
// relinking, use the empty string to leave this unspecified.
func (c *Client) UserSavedTracksOpt(limit int, offset uint, market string)(SavedTrackPage, error) {
  var result SavedTrackPage
  if limit > 50 || limit < 1 {
    return result, errors.New("spotify: SavedTracks supports a limit of 1 to 50")
  }
  spotifyURL := fmt.Sprintf("%sme/tracks?limit=%d", c.baseURL, limit)
  if offset != 0 {
    spotifyURL = fmt.Sprintf("%s&offset=%d", spotifyURL, offset)
  }
  if market != "" {
    spotifyURL = fmt.Sprintf("%s&market=%s", spotifyURL, market)
  }
  err := c.get(spotifyURL, &result)
  return result, err
}
