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
However, authenticated users benefit from increased rate limits,
even for calls that don't require authentication.

All functions requiring authentication are explicitly marked as
such in the godoc.  If you attempt to call any of these functions
without authenticating first, the error returned will be
`spotify.ErrNotAuthenticated`.

To authenticate [TBD ...]

## Examples

### Search

### Users

### Tracks

### Playlists
