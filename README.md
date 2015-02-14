Spotify
=======

[![GoDoc](https://godoc.org/github.com/zmb3/spotify?status.svg)](http://godoc.org/github.com/zmb3/spotify)

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
such in the godoc.  If you attempt to call any of these functions
without authorization, the error returned will be
`spotify.ErrAuthorizationRequired`.

Spotify uses OAuth2 for authentication, which typically requires the user to login
via a web browser.  Your application will have to implement the OAuth2 process
and provide the access token to this package.

The first step towards authenticating is to get a __client ID__ and __secret key__
by registering your application at the following page:

https://developer.spotify.com/my-applications/.

Use the ID and key to get an OAuth2 access token, and provide it to the client.
Additionally, specify the token type (most likely "Bearer").

````Go
c := spotify.Client{}
c.AccessToken = "my_token"
c.TokenType = spotify.BearerToken
````

You may find the following resources helpful:

1. Spotify's Web API Authorization Guide:
https://developer.spotify.com/web-api/authorization-guide/

2. Go's OAuth2 package:
https://godoc.org/golang.org/x/oauth2/google

3. spoticli - an example application that authenticates with OAuth2
https://github.com/zmb3/spoticli

## Helpful Hints

### Default Client

In general, for API calls that require authorization, you should create your own
`spotify.Client`.  For calls that don't require authorization, package level functions
are provided.  (These functions just proxy through `spotify.DefaultClient`, similar
to the way the `net/http` package works.)

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
    for _, item := results.Albums.Albums {
        fmt.Println("Album: ", item.Name)
    }
}
// handle playlist results
if results.Playlists != nil {
    for _, item := results.Playlists.Playlists {
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

To get information about the current user, you must authenticcate first:

````Go
c := spotify.Client{}
c.AccessToken = "my-token"
c.TokenType = spotify.BearerToken
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
