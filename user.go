package spotify

import (
	"encoding/json"
	"net/http"
)

// User contains the basic, publicly available
// information about a Spotify user.
type User struct {
	// The name displayed on the user's profile.
	// Note: Spotify currently fails to populate
	// this field when querying for a playlist.
	DisplayName string `json:"display_name"`
	// Known public external URLs for the user.
	ExternalURLs ExternalURL `json:"external_urls"`
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

// PrivateUser contains additional information about
// a user.  This data is private and requires user
// authentication.
type PrivateUser struct {
	User
	// the country of the user, as set in the user's
	// account profile.  An ISO 3166-1 alpha-2 country
	// code.  This field is only available when the
	// current user has granted acess to the
	// user-read-private scope.
	Country string `json:"country"`
	// The user's email address, as entered by the user
	// when creating their account.  Note: this email
	// is UNVERIFIED - there is no proof that it actually
	// belongs to the user.  This field is only available
	// when the current user has granted access to the
	// user-read-email scope.
	Email string `json:"email"`
	// The user's Spotify subscription level:
	// "premium", "free", etc.  The subscription level
	// "open" can be considered the same as "free".
	// This field is only available when the current user
	// has granted access to the user-read-private scope.
	Product string `json:"product"`
}

// UserPublicProfile is a wrapper around DefaultClient.UserPublicProfile.
func UserPublicProfile(userID ID) (*User, error) {
	return DefaultClient.UserPublicProfile(userID)
}

// UserPublicProfile gets public profile information about a
// Spotify User.  It does not require authentication.
func (c *Client) UserPublicProfile(userID ID) (*User, error) {
	uri := baseAddress + "users/" + string(userID)
	resp, err := c.http.Get(uri)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var spotifyError struct {
			E Error `json:"error"`
		}
		err = json.NewDecoder(resp.Body).Decode(&spotifyError)
		if err != nil {
			return nil, err
		}
		return nil, spotifyError.E
	}

	var user User
	err = json.NewDecoder(resp.Body).Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// CurrentUser gets detailed profile information about the
// current user.  It requires that the user authenticates first.
func (c *Client) CurrentUser() (*PrivateUser, error) {
	// TODO:
	//uri := baseAddress + "me"
	panic("spotify: CurrentUser() requires authentication!")
}
