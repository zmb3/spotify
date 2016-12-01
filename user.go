package spotify

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

// User contains the basic, publicly available information about a Spotify user.
type User struct {
	// The name displayed on the user's profile.
	// Note: Spotify currently fails to populate
	// this field when querying for a playlist.
	DisplayName string `json:"display_name"`
	// Known public external URLs for the user.
	ExternalURLs map[string]string `json:"external_urls"`
	// Information about followers of the user.
	Followers Followers `json:"followers"`
	// A link to the Web API endpoint for this user.
	Endpoint string `json:"href"`
	// The Spotify user ID for the user.
	ID string `json:"id"`
	// The user's profile image.
	Images []Image `json:"images"`
	// The Spotify URI for the user.
	URI URI `json:"uri"`
}

// PrivateUser contains additional information about a user.
// This data is private and requires user authentication.
type PrivateUser struct {
	User
	// The country of the user, as set in the user's account profile.
	// An ISO 3166-1 alpha-2 country code.  This field is only available when the
	// current user has granted acess to the ScopeUserReadPrivate scope.
	Country string `json:"country"`
	// The user's email address, as entered by the user when creating their account.
	// Note: this email is UNVERIFIED - there is no proof that it actually
	// belongs to the user.  This field is only available when the current user
	// has granted access to the ScopeUserReadEmail scope.
	Email string `json:"email"`
	// The user's Spotify subscription level: "premium", "free", etc.
	// The subscription level "open" can be considered the same as "free".
	// This field is only available when the current user has granted access to
	// the ScopeUserReadPrivate scope.
	Product string `json:"product"`
	// The user's date of birth, in the format 'YYYY-MM-DD'.  You can use
	// the DateLayout constant to convert this to a time.Time value.
	// This field is only available when the current user has granted
	// access to the ScopeUserReadBirthdate scope.
	Birthdate string `json:"birthdate"`
}

// GetUsersPublicProfile is a wrapper around DefaultClient.GetUsersPublicProfile.
func GetUsersPublicProfile(userID ID) (*User, error) {
	return DefaultClient.GetUsersPublicProfile(userID)
}

// GetUsersPublicProfile gets public profile information about a
// Spotify User.  It does not require authentication.
func (c *Client) GetUsersPublicProfile(userID ID) (*User, error) {
	spotifyURL := baseAddress + "users/" + string(userID)
	resp, err := c.http.Get(spotifyURL)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, decodeError(resp.Body)
	}
	var user User
	err = json.NewDecoder(resp.Body).Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// CurrentUser gets detailed profile information about the
// current user.  This call requires authorization.
//
// Reading the user's email address requires that the application
// has the ScopeUserReadEmail scope.  Reading the country, display
// name, profile images, and product subscription level requires
// that the application has the ScopeUserReadPrivate scope.
//
// Warning: The email address in the response will be the address
// that was entered when the user created their spotify account.
// This email address is unverified - do not assume that Spotify has
// checked that the email address actually belongs to the user.
func (c *Client) CurrentUser() (*PrivateUser, error) {
	resp, err := c.http.Get(baseAddress + "me")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, decodeError(resp.Body)
	}
	var result PrivateUser
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// CurrentUsersTracks gets a list of songs saved in the current
// Spotify user's "Your Music" library.
func (c *Client) CurrentUsersTracks() (*SavedTrackPage, error) {
	return c.CurrentUsersTracksOpt(nil)
}

// CurrentUsersTracksOpt is like CurrentUsersTracks, but it accepts additional
// options for sorting and filtering the results.
func (c *Client) CurrentUsersTracksOpt(opt *Options) (*SavedTrackPage, error) {
	spotifyURL := baseAddress + "me/tracks"
	if opt != nil {
		v := url.Values{}
		if opt.Country != nil {
			v.Set("country", *opt.Country)
		}
		if opt.Limit != nil {
			v.Set("limit", strconv.Itoa(*opt.Limit))
		}
		if opt.Offset != nil {
			v.Set("offset", strconv.Itoa(*opt.Offset))
		}
		if params := v.Encode(); params != "" {
			spotifyURL += "?" + params
		}
	}
	resp, err := c.http.Get(spotifyURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, decodeError(resp.Body)
	}
	var result SavedTrackPage
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// FollowUser adds the current user as a follower of one or more
// spotify users, identified by their Spotify IDs.
// This call requires authorization.
//
// Modifying the lists of artists or users the current user follows
// requires that the application has the ScopeUserFollowModify scope.
func (c *Client) FollowUser(ids ...ID) error {
	return c.modifyFollowers("user", true, ids...)
}

// FollowArtist adds the current user as a follower of one or more
// spotify artists, identified by their Spotify IDs.
// This call requires authorization.
//
// Modifying the lists of artists or users the current user follows
// requires that the application has the ScopeUserFollowModify scope.
func (c *Client) FollowArtist(ids ...ID) error {
	return c.modifyFollowers("artist", true, ids...)
}

// UnfollowUser removes the current user as a follower of one or more
// Spotify users.  This call requires authorization.
//
// Modifying the lists of artists or users the current user follows
// requires that the application has the ScopeUserFollowModify scope.
func (c *Client) UnfollowUser(ids ...ID) error {
	return c.modifyFollowers("user", false, ids...)
}

// UnfollowArtist removes the current user as a follower of one or more
// Spotify artists.  This call requires authorization.
//
// Modifying the lists of artists or users the current user follows
// requires that the application has the ScopeUserFollowModify scope.
func (c *Client) UnfollowArtist(ids ...ID) error {
	return c.modifyFollowers("artist", false, ids...)
}

// CurrentUserFollows checks to see if the current user is following
// one or more artists or other Spotify Users.  This call requires
// authorization, and that the application has the ScopeUserFollowRead
// scope.
//
// The t argument indicates the type of the IDs, and must be either
// "user" or "artist".
//
// The result is returned as a slice of bool values in the same order
// in which the IDs were specified.
func (c *Client) CurrentUserFollows(t string, ids ...ID) ([]bool, error) {
	if l := len(ids); l == 0 || l > 50 {
		return nil, errors.New("spotify: UserFollows supports 1 to 50 IDs")
	}
	if t != "artist" && t != "user" {
		return nil, errors.New("spotify: t must be 'artist' or 'user'")
	}
	spotifyURL := fmt.Sprintf("%sme/following/contains?type=%s&ids=%s",
		baseAddress, t, strings.Join(toStringSlice(ids), ","))
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
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (c *Client) modifyFollowers(usertype string, follow bool, ids ...ID) error {
	if l := len(ids); l == 0 || l > 50 {
		return errors.New("spotify: Follow/Unfollow supports 1 to 50 IDs")
	}
	v := url.Values{}
	v.Add("type", usertype)
	v.Add("ids", strings.Join(toStringSlice(ids), ","))
	spotifyURL := baseAddress + "me/following?" + v.Encode()
	method := "PUT"
	if !follow {
		method = "DELETE"
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
	if resp.StatusCode != http.StatusNoContent {
		return decodeError(resp.Body)
	}
	return nil
}

// CurrentUsersFollowedArtists gets the current user's followed artists.
// This call requires authorization, and that the user has granted the
// ScopeUserFollowRead scope.
func (c *Client) CurrentUsersFollowedArtists() (*FullArtistCursorPage, error) {
	return c.CurrentUsersFollowedArtistsOpt(-1, "")
}

// CurrentUsersFollowedArtistsOpt is like CurrentUsersFollowedArtists,
// but it accept the optional arguments limit and after.  Limit is the
// maximum number of items to return (1 <= limit <= 50), and after is
// the last artist ID retrieved from the previous request.  If you don't
// wish to specify either of the parameters, use -1 for limit and the empty
// string for after.
func (c *Client) CurrentUsersFollowedArtistsOpt(limit int, after string) (*FullArtistCursorPage, error) {
	spotifyURL := baseAddress + "me/following?type=artist"
	v := url.Values{}
	if limit != -1 {
		v.Set("limit", strconv.Itoa(limit))
	}
	if after != "" {
		v.Set("after", after)
	}
	if params := v.Encode(); params != "" {
		spotifyURL += "?" + params
	}
	resp, err := c.http.Get(spotifyURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, decodeError(resp.Body)
	}
	var result struct {
		A FullArtistCursorPage `json:"artists"`
	}
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return nil, err
	}
	return &result.A, nil
}

// CurrentUsersAlbums gets a list of albums saved in the current
// Spotify user's "Your Music" library.
func (c *Client) CurrentUsersAlbums() (*SavedAlbumPage, error) {
	return c.CurrentUsersAlbumsOpt(nil)
}

// CurrentUsersAlbumsOpt is like CurrentUsersAlbums, but it accepts additional
// options for sorting and filtering the results.
func (c *Client) CurrentUsersAlbumsOpt(opt *Options) (*SavedAlbumPage, error) {
	spotifyURL := baseAddress + "me/albums"
	if opt != nil {
		v := url.Values{}
		if opt.Country != nil {
			v.Set("market", *opt.Country)
		}
		if opt.Limit != nil {
			v.Set("limit", strconv.Itoa(*opt.Limit))
		}
		if opt.Offset != nil {
			v.Set("offset", strconv.Itoa(*opt.Offset))
		}
		if params := v.Encode(); params != "" {
			spotifyURL += "?" + params
		}
	}
	resp, err := c.http.Get(spotifyURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, decodeError(resp.Body)
	}
	var result SavedAlbumPage
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// CurrentUsersPlaylists gets a list of the playlists owned or followed by
// the current spotify user.  This call requires authorization.
// Private playlists require the ScopePlaylistReadPrivate scope.  Note that
// this scope alone will not return collaborative playlists, even though
// they are always private.  In order to retrieve collaborative playlists
// the user must authorize the ScopePlaylistReadCollaborative scope.
func (c *Client) CurrentUsersPlaylists() (*SimplePlaylistPage, error) {
	return c.CurrentUsersPlaylistsOpt(nil)
}

// CurrentUsersPlaylistsOpt is like CurrentUsersPlaylists, but it accepts
// additional options for sorting and filtering the results.
func (c *Client) CurrentUsersPlaylistsOpt(opt *Options) (*SimplePlaylistPage, error) {
	spotifyURL := baseAddress + "me/playlists"
	if opt != nil {
		v := url.Values{}
		if opt.Limit != nil {
			v.Set("limit", strconv.Itoa(*opt.Limit))
		}
		if opt.Offset != nil {
			v.Set("offset", strconv.Itoa(*opt.Offset))
		}
		if params := v.Encode(); params != "" {
			spotifyURL += "?" + params
		}
	}
	resp, err := c.http.Get(spotifyURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, decodeError(resp.Body)
	}
	var result SimplePlaylistPage
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}
