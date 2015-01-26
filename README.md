Spotify
=======

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

__ Authentication is a work in progress... __

Most of the functionality is available without authenticating.
However, authenticated users benefit from increased rate limits.

Features that access a user's private data require authorization.
All functions requiring authorization are explicitly marked as
such in the godoc.  If you attempt to call any of these functions
without authorization, the error returned will be
`spotify.ErrAuthorizationRequired`.

Spotify uses OAuth2 for authentication.

The first step towards authenticating is to get a __client ID__ and __secret key__
by registering your application at the following page:

https://developer.spotify.com/my-applications/.


## Examples

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

The search query supports a variety of operations.  
Refer to the godoc for more information.

### Users

### Tracks

### Playlists
