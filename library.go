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
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
)

// UserHasTracks checks if one or more tracks are saved to the current user's
// "Your Music" library.  This call requires authorization.
func (c *Client) UserHasTracks(ids ...ID) ([]bool, error) {
	if l := len(ids); l == 0 || l > 50 {
		return nil, errors.New("spotify: UserHasTracks supports 1 to 50 IDs per call")
	}
	spotifyURL := fmt.Sprintf("%sme/tracks/contains?ids=%s", baseAddress, strings.Join(toStringSlice(ids), ","))
	resp, err := c.http.Get(spotifyURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, decodeError(resp.Body)
	}
	var result []bool
	err = json.NewDecoder(resp.Body).Decode(&result)
	return result, err
}

// AddTracksToLibrary saves one or more tracks to the current user's
// "Your Music" library.  This call requires authorization (the
// ScopeUserLibraryModify scope).
// A track can only be saved once; duplicate IDs are ignored.
func (c *Client) AddTracksToLibrary(ids ...ID) error {
	return c.modifyLibraryTracks(true, ids...)
}

// RemoveTracksFromLibrary removes one or more tracks from the current user's
// "Your Music" library.  This call requires authorization (the ScopeUserModifyLibrary
// scope).  Trying to remove a track when you do not have the user's authorization
// results in a `spotify.Error` with the status code set to http.StatusUnauthorized.
func (c *Client) RemoveTracksFromLibrary(ids ...ID) error {
	return c.modifyLibraryTracks(false, ids...)
}

func (c *Client) modifyLibraryTracks(add bool, ids ...ID) error {
	if l := len(ids); l == 0 || l > 50 {
		return errors.New("spotify: this call supports 1 to 50 IDs per call")
	}
	spotifyURL := fmt.Sprintf("%sme/tracks?ids=%s", baseAddress, strings.Join(toStringSlice(ids), ","))
	method := "DELETE"
	if add {
		method = "PUT"
	}
	req, err := http.NewRequest(method, spotifyURL, nil)
	if err != nil {
		return err
	}
	resp, err := c.http.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return decodeError(resp.Body)
	}
	return nil
}
