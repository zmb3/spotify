package spotify

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
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

// GetUsersPublicProfile gets public profile information about a
// Spotify User.  It does not require authentication.
func (c *Client) GetUsersPublicProfile(ctx context.Context, userID ID) (*User, error) {
	spotifyURL := c.baseURL + "users/" + string(userID)

	var user User

	err := c.get(ctx, spotifyURL, &user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

// CurrentUser gets detailed profile information about the
// current user.
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
func (c *Client) CurrentUser(ctx context.Context) (*PrivateUser, error) {
	var result PrivateUser

	err := c.get(ctx, c.baseURL+"me", &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

// CurrentUsersShows gets a list of shows saved in the current
// Spotify user's "Your Music" library.
//
// API Doc: https://developer.spotify.com/documentation/web-api/reference-beta/#endpoint-get-users-saved-shows
//
// Supported options: Limit, Offset
func (c *Client) CurrentUsersShows(ctx context.Context, opts ...RequestOption) (*SavedShowPage, error) {
	spotifyURL := c.baseURL + "me/shows"
	if params := processOptions(opts...).urlParams.Encode(); params != "" {
		spotifyURL += "?" + params
	}

	var result SavedShowPage

	err := c.get(ctx, spotifyURL, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

// CurrentUsersTracks gets a list of songs saved in the current
// Spotify user's "Your Music" library.
//
// API Doc: https://developer.spotify.com/documentation/web-api/reference-beta/#endpoint-get-users-saved-tracks
//
// Supported options: Limit, Country, Offset
func (c *Client) CurrentUsersTracks(ctx context.Context, opts ...RequestOption) (*SavedTrackPage, error) {
	spotifyURL := c.baseURL + "me/tracks"
	if params := processOptions(opts...).urlParams.Encode(); params != "" {
		spotifyURL += "?" + params
	}

	var result SavedTrackPage

	err := c.get(ctx, spotifyURL, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

// FollowUser adds the current user as a follower of one or more
// spotify users, identified by their Spotify IDs.
//
// Modifying the lists of artists or users the current user follows
// requires that the application has the ScopeUserFollowModify scope.
func (c *Client) FollowUser(ctx context.Context, ids ...ID) error {
	return c.modifyFollowers(ctx, "user", true, ids...)
}

// FollowArtist adds the current user as a follower of one or more
// spotify artists, identified by their Spotify IDs.
//
// Modifying the lists of artists or users the current user follows
// requires that the application has the ScopeUserFollowModify scope.
func (c *Client) FollowArtist(ctx context.Context, ids ...ID) error {
	return c.modifyFollowers(ctx, "artist", true, ids...)
}

// UnfollowUser removes the current user as a follower of one or more
// Spotify users.
//
// Modifying the lists of artists or users the current user follows
// requires that the application has the ScopeUserFollowModify scope.
func (c *Client) UnfollowUser(ctx context.Context, ids ...ID) error {
	return c.modifyFollowers(ctx, "user", false, ids...)
}

// UnfollowArtist removes the current user as a follower of one or more
// Spotify artists.
//
// Modifying the lists of artists or users the current user follows
// requires that the application has the ScopeUserFollowModify scope.
func (c *Client) UnfollowArtist(ctx context.Context, ids ...ID) error {
	return c.modifyFollowers(ctx, "artist", false, ids...)
}

// CurrentUserFollows checks to see if the current user is following
// one or more artists or other Spotify Users.  This call requires
// ScopeUserFollowRead.
//
// The t argument indicates the type of the IDs, and must be either
// "user" or "artist".
//
// The result is returned as a slice of bool values in the same order
// in which the IDs were specified.
func (c *Client) CurrentUserFollows(ctx context.Context, t string, ids ...ID) ([]bool, error) {
	if l := len(ids); l == 0 || l > 50 {
		return nil, errors.New("spotify: UserFollows supports 1 to 50 IDs")
	}
	if t != "artist" && t != "user" {
		return nil, errors.New("spotify: t must be 'artist' or 'user'")
	}
	spotifyURL := fmt.Sprintf("%sme/following/contains?type=%s&ids=%s",
		c.baseURL, t, strings.Join(toStringSlice(ids), ","))

	var result []bool

	err := c.get(ctx, spotifyURL, &result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (c *Client) modifyFollowers(ctx context.Context, usertype string, follow bool, ids ...ID) error {
	if l := len(ids); l == 0 || l > 50 {
		return errors.New("spotify: Follow/Unfollow supports 1 to 50 IDs")
	}
	v := url.Values{}
	v.Add("type", usertype)
	v.Add("ids", strings.Join(toStringSlice(ids), ","))
	spotifyURL := c.baseURL + "me/following?" + v.Encode()
	method := "PUT"
	if !follow {
		method = "DELETE"
	}
	req, err := http.NewRequestWithContext(ctx, method, spotifyURL, nil)
	if err != nil {
		return err
	}
	err = c.execute(req, nil, http.StatusNoContent)
	if err != nil {
		return err
	}
	return nil
}

// CurrentUsersFollowedArtists gets the current user's followed artists.
// This call requires that the user has granted the ScopeUserFollowRead scope.
// Supported options: Limit, After
func (c *Client) CurrentUsersFollowedArtists(ctx context.Context, opts ...RequestOption) (*FullArtistCursorPage, error) {
	spotifyURL := c.baseURL + "me/following"
	v := processOptions(opts...).urlParams
	v.Set("type", "artist")
	if params := v.Encode(); params != "" {
		spotifyURL += "?" + params
	}

	var result struct {
		A FullArtistCursorPage `json:"artists"`
	}

	err := c.get(ctx, spotifyURL, &result)
	if err != nil {
		return nil, err
	}

	return &result.A, nil
}

// CurrentUsersAlbums gets a list of albums saved in the current
// Spotify user's "Your Music" library.
//
// Supported options: Market, Limit, Offset
func (c *Client) CurrentUsersAlbums(ctx context.Context, opts ...RequestOption) (*SavedAlbumPage, error) {
	spotifyURL := c.baseURL + "me/albums"
	if params := processOptions(opts...).urlParams.Encode(); params != "" {
		spotifyURL += "?" + params
	}

	var result SavedAlbumPage

	err := c.get(ctx, spotifyURL, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

// CurrentUsersPlaylists gets a list of the playlists owned or followed by
// the current spotify user.
//
// Private playlists require the ScopePlaylistReadPrivate scope.  Note that
// this scope alone will not return collaborative playlists, even though
// they are always private.  In order to retrieve collaborative playlists
// the user must authorize the ScopePlaylistReadCollaborative scope.
//
// Supported options: Limit, Offset
func (c *Client) CurrentUsersPlaylists(ctx context.Context, opts ...RequestOption) (*SimplePlaylistPage, error) {
	spotifyURL := c.baseURL + "me/playlists"
	if params := processOptions(opts...).urlParams.Encode(); params != "" {
		spotifyURL += "?" + params
	}

	var result SimplePlaylistPage

	err := c.get(ctx, spotifyURL, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

// CurrentUsersTopArtists fetches a list of the user's top artists over the specified Timerange.
// The default is medium term.
//
// Supported options: Limit, Timerange
func (c *Client) CurrentUsersTopArtists(ctx context.Context, opts ...RequestOption) (*FullArtistPage, error) {
	spotifyURL := c.baseURL + "me/top/artists"
	if params := processOptions(opts...).urlParams.Encode(); params != "" {
		spotifyURL += "?" + params
	}

	var result FullArtistPage

	err := c.get(ctx, spotifyURL, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

// CurrentUsersTopTracks fetches the user's top tracks over the specified Timerange
// sensible defaults. The default limit is 20 and the default timerange
// is medium_term. This call requires ScopeUserTopRead.
//
// Supported options: Limit, Timerange, Offset
func (c *Client) CurrentUsersTopTracks(ctx context.Context, opts ...RequestOption) (*FullTrackPage, error) {
	spotifyURL := c.baseURL + "me/top/tracks"
	if params := processOptions(opts...).urlParams.Encode(); params != "" {
		spotifyURL += "?" + params
	}

	var result FullTrackPage

	err := c.get(ctx, spotifyURL, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}
