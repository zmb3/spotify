package spotify

import (
	"context"
)

// User contains the basic, publicly available information about a Spotify user.
type User struct {
	// The name displayed on the user's profile.
	// Note: Spotify currently fails to populate
	// this field when querying for a playlist.
	DisplayName string `json:"display_name"`
	// Known public external URLs for the user.
	ExternalURLs map[string]string `json:"external_urls"`
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
	// The user's date of birth, in the format 'YYYY-MM-DD'.  You can use
	// [DateLayout] to convert this to a [time.Time] value. This field is only
	// available when the current user has granted access to the
	// [ScopeUserReadBirthdate] scope.
	Birthdate string `json:"birthdate"`
}

// CurrentUser gets detailed profile information about the
// [current user].
//
// Reading the user's email address requires that the application
// has the [ScopeUserReadEmail] scope.  Reading the country, display
// name, profile images, and product subscription level requires
// that the application has the [ScopeUserReadPrivate] scope.
//
// Warning: The email address in the response will be the address
// that was entered when the user created their spotify account.
// This email address is unverified - do not assume that Spotify has
// checked that the email address actually belongs to the user.
//
// [current user]: https://developer.spotify.com/documentation/web-api/reference/get-current-users-profile
func (c *Client) CurrentUser(ctx context.Context) (*PrivateUser, error) {
	var result PrivateUser

	err := c.get(ctx, c.baseURL+"me", &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

// CurrentUsersShows gets a [list of shows] saved in the current
// Spotify user's "Your Music" library.
//
// Supported options: [Limit], [Offset].
//
// [list of shows]: https://developer.spotify.com/documentation/web-api/reference/get-users-saved-shows
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

// CurrentUsersTracks gets a [list of songs] saved in the current
// Spotify user's "Your Music" library.
//
// Supported options: [Limit], [Country], [Offset].
//
// [list of songs]: https://developer.spotify.com/documentation/web-api/reference/get-users-saved-tracks
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

// CurrentUsersFollowedArtists gets the [current user's followed artists].
// This call requires that the user has granted the [ScopeUserFollowRead] scope.
//
// Supported options: [Limit], [After].
//
// [current user's followed artists]: https://developer.spotify.com/documentation/web-api/reference/get-followed
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

// CurrentUsersAlbums gets a [list of albums] saved in the current
// Spotify user's "Your Music" library.
//
// Supported options: [Market], [Limit], [Offset].
//
// [list of albums]: https://developer.spotify.com/documentation/web-api/reference/get-users-saved-albums
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

// CurrentUsersPlaylists gets a [list of the playlists] owned or followed by
// the current spotify user.
//
// Private playlists require the [ScopePlaylistReadPrivate] scope.  Note that
// this scope alone will not return collaborative playlists, even though
// they are always private.  In order to retrieve collaborative playlists
// the user must authorize the [ScopePlaylistReadCollaborative] scope.
//
// Supported options: [Limit], [Offset].
//
// [list of the playlists]: https://developer.spotify.com/documentation/web-api/reference/get-a-list-of-current-users-playlists
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

// CurrentUsersTopArtists fetches a list of the [user's top artists] over the specified [Timerange].
// The default is [MediumTermRange].
//
// Supported options: [Limit], [Timerange]
//
// [user's top artists]: https://developer.spotify.com/documentation/web-api/reference/get-users-top-artists-and-tracks
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

// CurrentUsersTopTracks fetches the [user's top tracks] over the specified
// [Timerange]. The default limit is 20 and the default timerange is
// [MediumTermRange]. This call requires [ScopeUserTopRead].
//
// Supported options: [Limit], [Timerange], [Offset].
//
// [user's top tracks]: https://developer.spotify.com/documentation/web-api/reference/get-users-top-artists-and-tracks
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
