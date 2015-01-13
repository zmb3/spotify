Spotify
=======

This is a Go wrapper for working with Spotify's
[Web API](https://developer.spotify.com/web-api/).

It supports every task listed in the Web API Endpoint Reference,
located [here](https://developer.spotify.com/web-api/endpoint-reference/).

By using this library you agree to Spotify's
[Developer Terms of Use](https://developer.spotify.com/developer-terms-of-use/).

## Installation

To install the library, simply

`go get github.com/zmb3/spotify`

## Authentication

Most of the functionality is available without authenticating.
However, authenticated users benefit from increased rate limits.

Features that access a user's private data require authorization.
All functions requiring authorization are explicitly marked as
such in the godoc.  If you attempt to call any of these functions
without authorization, the error returned will be
`spotify.ErrNotAuthorized`.

The first step towards authenticating is to get a __client ID__ and __secret key__ by registering your application at the following page:

https://developer.spotify.com/my-applications/.


## Examples

### Search

### Users

### Tracks

### Playlists
