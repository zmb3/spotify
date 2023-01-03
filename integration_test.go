package spotify

import (
	"context"
	"os"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	spotifyauth "github.com/zmb3/spotify/v2/auth"
	"golang.org/x/oauth2/clientcredentials"
)

var (
	integrationSetupOnce = sync.Once{}
	integrationSetupErr  error
	// integrationSetupResult MUST be accessed via integrationTest.
	integrationSetupResult *Client
)

func integrationTest(t *testing.T) *Client {
	t.Helper()
	if testing.Short() {
		t.Skip("Flag -short provided. Skipping integration test.")
	}

	integrationSetupOnce.Do(func() {
		ctx := context.Background()
		config := &clientcredentials.Config{
			ClientID:     os.Getenv("SPOTIFY_ID"),
			ClientSecret: os.Getenv("SPOTIFY_SECRET"),
			TokenURL:     spotifyauth.TokenURL,
		}
		token, err := config.Token(ctx)
		if err != nil {
			integrationSetupErr = err
			return
		}
		httpClient := spotifyauth.New().Client(ctx, token)
		integrationSetupResult = New(httpClient)
	})
	require.NoError(t, integrationSetupErr)

	return integrationSetupResult
}

func TestClient_GetTrack_Integration(t *testing.T) {
	// Black Country, New Road - Sunglasses
	// https://open.spotify.com/track/1sT5Wh3SVv6nhs7lgPEnkl
	t.Parallel()
	c := integrationTest(t)
	ctx := context.Background()

	track, err := c.GetTrack(ctx, ID("1sT5Wh3SVv6nhs7lgPEnkl"))
	require.NoError(t, err)

	artist := SimpleArtist{
		Name:     "Black Country, New Road",
		ID:       "3PP6ghmOlDl2jaKaH0avUN",
		URI:      "spotify:artist:3PP6ghmOlDl2jaKaH0avUN",
		Endpoint: "https://api.spotify.com/v1/artists/3PP6ghmOlDl2jaKaH0avUN",
		ExternalURLs: map[string]string{
			"spotify": "https://open.spotify.com/artist/3PP6ghmOlDl2jaKaH0avUN",
		},
	}
	// SimpleTrack
	assert.Equal(t, []SimpleArtist{artist}, track.Artists)
	// omit tight check on available markets as this value fluctuates too
	// often.
	assert.NotEmpty(t, track.AvailableMarkets)
	assert.Equal(t, 590753, track.Duration)
	assert.Equal(t, map[string]string{
		"spotify": "https://open.spotify.com/track/1sT5Wh3SVv6nhs7lgPEnkl",
	}, track.ExternalURLs)
	assert.Equal(t, "https://api.spotify.com/v1/tracks/1sT5Wh3SVv6nhs7lgPEnkl", track.Endpoint)
	assert.Equal(t, "1sT5Wh3SVv6nhs7lgPEnkl", track.ID)
	assert.Equal(t, "Sunglasses", track.Name)
	assert.Equal(t, "https://p.scdn.co/mp3-preview/f30122eb0fa9408f796107f43b869292d004a42a?cid=4ea3eb5b8ba541319c23c51a508b6b56", track.PreviewURL)
	assert.Equal(t, 4, track.TrackNumber)
	assert.Equal(t, "spotify:track:1sT5Wh3SVv6nhs7lgPEnkl", track.URI)
	assert.Equal(t, "track", track.Type)
	// SimpleAlbum
	assert.Equal(t, []SimpleArtist{artist}, track.Album.Artists)
}
