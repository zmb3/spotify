# Contributing

## Guidance

- PRs must include tests.
- New tests should aim to leverage `testify/assert` and `testify/require`

## Running integration tests

Create an application in the Spotify developer console.

Store your client ID and secret in SPOTIFY_ID and SPOTIFY_SECRET environment
variables.

```sh
INTEGRATION_TEST=1 go test ./...
```