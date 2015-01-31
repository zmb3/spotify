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
	"net/http"
	"net/url"
	"os"
	"strings"

	"golang.org/x/oauth2"
)

// TokenType indicates which type of authorization a client uses.
type TokenType string

const (
	// BasicToken authorization provides an increased rate limit, but don't
	// offer access to a user's private data.
	BasicToken TokenType = "Basic"
	// BearerToken authorization offers access to a user's private data.
	BearerToken = "Bearer"
)

const (
	// AuthorizeBaseAddress is the URL to the Spotify Accounts Service's
	// authorization endpoint.
	AuthorizeBaseAddress = "https://accounts.spotify.com/authorize"
	// TokenBaseAddress is the URL to the Spotify Accounts Service's
	// token endpoint.
	TokenBaseAddress = "https://accounts.spotify.com/api/token"
)

// Oauth2Endpoint contains the OAuth2 token endpoint URLs for the Spotify Web API.
var Oauth2Endpoint = oauth2.Endpoint{
	AuthURL:  AuthorizeBaseAddress,
	TokenURL: TokenBaseAddress,
}

// Scopes let you specify exactly which types of data your
// application wants to access.  The set of scopes you pass
// in your authentication request determines what access the
// permissions the user is asked to grant.
const (
	// ScopePlaylistReadPrivate seeks permission to read
	// a user's private playlists.
	ScopePlaylistReadPrivate = "playlist-read-private"
	// ScopePlaylistModifyPublic seeks write access
	// to a user's public playlists.
	ScopePlaylistModifyPublic = "playlist-modify-public"
	// ScopePlaylistModifyPrivate seeks write access to
	// a user's private playlists.
	ScopePlaylistModifyPrivate = "playlist-modify-private"
	// ScopeUserFollowModify seeks write/delete access to
	// the list of artists and other users that a user follows.
	ScopeUserFollowModify = "user-follow-modify"
	// ScopeUserFollowRead seeks read access to the list of
	// artists and other users that a user follows.
	ScopeUserFollowRead = "user-follow-read"
	// ScopeUserLibraryModify seeks write/delete acess to a
	// user's "Your Music" library.
	ScopeUserLibraryModify = "user-library-modify"
	// ScopeUserLibraryRead seeks read access to a user's
	// "Your Music" library.
	ScopeUserLibraryRead = "user-library-read"
	// ScopeUserReadPrivate seeks read access to a user's
	// subsription details (type of user account)
	ScopeUserReadPrivate = "user-read-private"
	// ScopeUserReadEmail seeks read access to a user's
	// email address.
	ScopeUserReadEmail = "user-read-email"
)

// AuthenticationOptions contains paramters required for
// the client credentials authentication flow.
type AuthenticationOptions struct {
	// Scopes is a list of scopes that specify exactly which
	// types of information your application plans to access.
	// If nil, then authorization will be granted only to access
	// publicly available information.
	Scopes []string
	// Client ID identifies your application.  Get one by registering
	// at https://developer.spotify.com/my-applications/.
	// If nil, then the client ID will be read from the SPOTIFY_CLIENT_ID
	// environment variable.
	ClientID *string
	// ClientSecret is the secret key for your application.  If nil, then
	// the secret key will be read from the SPOTIFY_SECRET_KEY
	// environment variable.
	ClientSecret *string
}

// AuthenticateClientCredentials uses the client credentials flow,
// which makes it possible to authenticate your requests to the
// Spotify Web API in order to obtain a higher rate limit.  This
// flow does NOT include authorization to access a user's private data.
func (c *Client) AuthenticateClientCredentials(opt AuthenticationOptions) error {
	var id, secret string
	if opt.ClientID == nil {
		id = os.Getenv("SPOTIFY_ID")
	} else {
		id = *opt.ClientID
	}
	if opt.ClientSecret == nil {
		secret = os.Getenv("SPOTIFY_SECRET")
	} else {
		secret = *opt.ClientSecret
	}

	if id == "" || secret == "" {
		return errors.New("spotify: missing client ID/secret key")
	}
	values := url.Values{}
	values.Set("grant_type", "client_credentials")

	if opt.Scopes != nil {
		values.Set("scopes", strings.Join(opt.Scopes, " "))
	}

	req, err := http.NewRequest("POST", TokenBaseAddress+"?"+values.Encode(), nil)
	if err != nil {
		return err
	}
	req.SetBasicAuth(id, secret)
	resp, err := c.http.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	var body struct {
		AccessToken string `json:"access_token"`
		ExpiresIn   int    `json:"expires_in"`
	}
	err = json.NewDecoder(resp.Body).Decode(&body)
	if err != nil {
		return err
	}
	// TODO: c.AccessToken = body.AccessToken
	// TODO: c.TokenExpiration = ...

	// now the client has a non-nil token
	// TODO: all api calls must be udpated to include access token in header
	return nil
}
