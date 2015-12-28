
Spotify
=======

[![GoDoc](https://godoc.org/github.com/zmb3/spotify?status.svg)](http://godoc.org/github.com/zmb3/spotify)
[![Build status](https://ci.appveyor.com/api/projects/status/1nr9vv0jqq438nj2?svg=true)](https://ci.appveyor.com/project/zmb3/spotify)
[![Build Status](https://travis-ci.org/zmb3/spotify.svg)](https://travis-ci.org/zmb3/spotify)

This is a Go wrapper for working with Spotify's
[Web API](https://developer.spotify.com/web-api/).

It aims to support every task listed in the Web API Endpoint Reference,
located [here](https://developer.spotify.com/web-api/endpoint-reference/).

By using this library you agree to Spotify's
[Developer Terms of Use](https://developer.spotify.com/developer-terms-of-use/).

## Installation

To install the library, simply

`go get github.com/zmb3/spotify`

## Authentication

Most of the Web API functionality is available without authenticating.
However, authenticated users benefit from increased rate limits.

Features that access a user's private data require authorization.
All functions requiring authorization are explicitly marked as
such in the godoc.

Spotify uses OAuth2 for authentication, which typically requires the user to login
via a web browser.  This package includes an `Authenticator` type to handle the details for you.

Start by getting registering your application at the following page:

https://developer.spotify.com/my-applications/.

You'll get a __client ID__ and __secret key__ for your application.  An easy way to
provide this data to your application is to set the SPOTIFY_ID and SPOTIFY_SECRET
environment variables.  If you choose not to use environment variables, you can
provide this data manually.


````Go
// the redirect URL must be an exact match of a URL you've registered for your application
// scopes determine which permissions the user is prompted to authorize
auth := spotify.NewAuthenticator(redirectURL, spotify.ScopeUserReadPrivate)

// if you didn't store your ID and secret key in the specified environment variables,
// you can set them manually here
auth.SetAuthInfo(clientID, secretKey)

// get the user to this URL - how you do that is up to you
// you should specify a unique state string to identify the session
url := auth.AuthURL(state)

// the user will eventually be redirected back to your redirect URL
// typically you'll have a handler set up like the following:
func redirectHandler(w http.ResponseWriter, r *http.Request) {
      // use the same state string here that you used to generate the URL
      token, err := auth.Token(state, r)
      if err != nil {
            http.Error(w, "Couldn't get token", http.StatusNotFound)
            return
      }
      // create a client using the specified token
      client := auth.NewClient(token)

      // the client can now be used to make authenticated requests
}
````

You may find the following resources useful:

1. Spotify's Web API Authorization Guide:
https://developer.spotify.com/web-api/authorization-guide/

2. Go's OAuth2 package:
https://godoc.org/golang.org/x/oauth2/google


## Helpful Hints

### Default Client

For API calls that require authorization, you should create your own
`spotify.Client` using an `Authenticator`.  For calls that don't require authorization,
package level wrapper functions are provided (see `spotify.Search` for example)

These functions just proxy through `spotify.DefaultClient`, similar to the way
the `net/http` package works.

### Optional Parameters

Many of the functions in this package come in two forms - a simple version that
omits optional parameters and uses reasonable defaults, and a more sophisticated
version that accepts additional parameters.  The latter is suffixed with `Opt`
to indicate that it accepts some optional parameters.

## API Examples

### Search

The search functionality returns a set of results, grouped by type (album, artist,
playlist, and track).  You can search for more than one type of item with a single
search.  For example, to search for holiday playlists and albums:

```Go
results, err := spotify.Search("holiday", SearchTypePlaylist|SearchTypeAlbum)
// error handling omitted

// handle album results
if results.Albums != nil {
    for _, item := range results.Albums.Albums {
        fmt.Println("Album: ", item.Name)
    }
}
// handle playlist results
if results.Playlists != nil {
    for _, item := range results.Playlists.Playlists {
        fmt.Println("Playlist: ", item.Name)
    }
}
```

The search query supports a variety of queries.  Refer to the godoc for more information.

### Users

To get information about another Spotify user:

````Go
user, err := spotify.GetUsersPublicProfile("user-id")
// err handling omitted
fmt.Println(user.DisplayName)
fmt.Println(user.Followers.Count, "followers")
````

To get information about the current user, you must authenticate first:

````Go
c := spotify.Client{}
c.HTTP = getHTTPClient()
me, err := c.CurrentUser()
// error handling omitted

fmt.Println(me.ID, me.Email, me.DisplayName)
````

To check if a user follows another user:

````Go
// authentication omitted
follows, err := c.CurrentUserFollows("user", "<spotify_id_here")
// error handling omitted
if follows[0] {
    fmt.Println("Yes, the current user follows this user.")
} else {
    fmt.Println("No, the current user does not follow this user")
}
````

You can check multiple items in the same call (and this works for artists too):

````Go
// authentication omitted
follows, err := c.CurrentUserFollows("artist", "artist_0_id", "artist_1_id", "artist_2_id")
// error handling omitted
if follows[0] {
    fmt.Println("The user follows artist 0")
}
if follows[1] {
    fmt.Println("The user follows artist 1")
}
if follows[2] {
    fmt.Println("The user follows artist 2")
}

````

### Tracks

To get catalog information for a track:

````Go
track, err := spotify.GetTrack("track_id_here")
// error handling omitted
fmt.Printf("%s is %d milliseconds long and has a popularity of %d\n",
    track.Name, track.Duration, track.Popularity)
````

### Playlists

To get a list of Spotify's featured playlists, authenticate first, and then:

````Go
msg, page, err := c.FeaturedPlaylists()
// error handling omitted

for _, playlist := range page.Playlists {
    fmt.Println(playlist.Name, playlist.Owner)
}

````
