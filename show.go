package spotify

import (
	"context"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type SavedShow struct {
	// The date and time the show was saved, represented as an ISO
	// 8601 UTC timestamp with a zero offset (YYYY-MM-DDTHH:MM:SSZ).
	// You can use the TimestampLayout constant to convert this to
	// a time.Time value.
	AddedAt  string `json:"added_at"`
	FullShow `json:"show"`
}

// FullShow contains full data about a show.
type FullShow struct {
	SimpleShow

	// A list of the show’s episodes.
	Episodes SimpleEpisodePage `json:"episodes"`
}

// SimpleShow contains basic data about a show.
type SimpleShow struct {
	// A list of the countries in which the show can be played,
	// identified by their ISO 3166-1 alpha-2 code.
	AvailableMarkets []string `json:"available_markets"`

	// The copyright statements of the show.
	Copyrights []Copyright `json:"copyrights"`

	// A description of the show.
	Description string `json:"description"`

	// Whether or not the show has explicit content
	// (true = yes it does; false = no it does not OR unknown).
	Explicit bool `json:"explicit"`

	// Known external URLs for this show.
	ExternalURLs map[string]string `json:"external_urls"`

	// A link to the Web API endpoint providing full details
	// of the show.
	Href string `json:"href"`

	// The SpotifyID for the show.
	ID ID `json:"id"`

	// The cover art for the show in various sizes,
	// widest first.
	Images []Image `json:"images"`

	// True if all of the show’s episodes are hosted outside
	// of Spotify’s CDN. This field might be null in some cases.
	IsExternallyHosted *bool `json:"is_externally_hosted"`

	// A list of the languages used in the show, identified by
	// their ISO 639 code.
	Languages []string `json:"languages"`

	// The media type of the show.
	MediaType string `json:"media_type"`

	// The name of the show.
	Name string `json:"name"`

	// The publisher of the show.
	Publisher string `json:"publisher"`

	// The object type: “show”.
	Type string `json:"type"`

	// The Spotify URI for the show.
	URI URI `json:"uri"`
}

type EpisodePage struct {
	// A URL to a 30 second preview (MP3 format) of the episode.
	AudioPreviewURL string `json:"audio_preview_url"`

	// A description of the episode.
	Description string `json:"description"`

	// The episode length in milliseconds.
	Duration_ms Numeric `json:"duration_ms"`

	// Whether or not the episode has explicit content
	// (true = yes it does; false = no it does not OR unknown).
	Explicit bool `json:"explicit"`

	// 	External URLs for this episode.
	ExternalURLs map[string]string `json:"external_urls"`

	// A link to the Web API endpoint providing full details of the episode.
	Href string `json:"href"`

	// The Spotify ID for the episode.
	ID ID `json:"id"`

	// The cover art for the episode in various sizes, widest first.
	Images []Image `json:"images"`

	// True if the episode is hosted outside of Spotify’s CDN.
	IsExternallyHosted bool `json:"is_externally_hosted"`

	// True if the episode is playable in the given market.
	// Otherwise false.
	IsPlayable bool `json:"is_playable"`

	// A list of the languages used in the episode, identified by their ISO 639 code.
	Languages []string `json:"languages"`

	// The name of the episode.
	Name string `json:"name"`

	// The date the episode was first released, for example
	// "1981-12-15". Depending on the precision, it might
	// be shown as "1981" or "1981-12".
	ReleaseDate string `json:"release_date"`

	// The precision with which release_date value is known:
	// "year", "month", or "day".
	ReleaseDatePrecision string `json:"release_date_precision"`

	// The user’s most recent position in the episode. Set if the
	// supplied access token is a user token and has the scope
	// user-read-playback-position.
	ResumePoint ResumePointObject `json:"resume_point"`

	// The show on which the episode belongs.
	Show SimpleShow `json:"show"`

	// The object type: "episode".
	Type string `json:"type"`

	// The Spotify URI for the episode.
	URI URI `json:"uri"`
}

type ResumePointObject struct {
	// 	Whether or not the episode has been fully played by the user.
	FullyPlayed bool `json:"fully_played"`

	// The user’s most recent position in the episode in milliseconds.
	ResumePositionMs Numeric `json:"resume_position_ms"`
}

// ReleaseDateTime converts the show's ReleaseDate to a time.TimeValue.
// All of the fields in the result may not be valid.  For example, if
// ReleaseDatePrecision is "month", then only the month and year
// (but not the day) of the result are valid.
func (e *EpisodePage) ReleaseDateTime() time.Time {
	if e.ReleaseDatePrecision == "day" {
		result, _ := time.Parse(DateLayout, e.ReleaseDate)
		return result
	}
	if e.ReleaseDatePrecision == "month" {
		ym := strings.Split(e.ReleaseDate, "-")
		year, _ := strconv.Atoi(ym[0])
		month, _ := strconv.Atoi(ym[1])
		return time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
	}
	year, _ := strconv.Atoi(e.ReleaseDate)
	return time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC)
}

// GetShow retrieves information about a specific show.
// API reference: https://developer.spotify.com/documentation/web-api/reference/#endpoint-get-a-show
// Supported options: Market
func (c *Client) GetShow(ctx context.Context, id ID, opts ...RequestOption) (*FullShow, error) {
	spotifyURL := c.baseURL + "shows/" + string(id)
	if params := processOptions(opts...).urlParams.Encode(); params != "" {
		spotifyURL += "?" + params
	}

	var result FullShow

	err := c.get(ctx, spotifyURL, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

// GetShowEpisodes retrieves paginated episode information about a specific show.
// API reference: https://developer.spotify.com/documentation/web-api/reference/#endpoint-get-a-shows-episodes
// Supported options: Market, Limit, Offset
func (c *Client) GetShowEpisodes(ctx context.Context, id string, opts ...RequestOption) (*SimpleEpisodePage, error) {
	spotifyURL := c.baseURL + "shows/" + id + "/episodes"
	if params := processOptions(opts...).urlParams.Encode(); params != "" {
		spotifyURL += "?" + params
	}

	var result SimpleEpisodePage

	err := c.get(ctx, spotifyURL, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

// SaveShowsForCurrentUser saves one or more shows to current Spotify user's library.
// API reference: https://developer.spotify.com/documentation/web-api/reference/#/operations/save-shows-user
func (c *Client) SaveShowsForCurrentUser(ctx context.Context, ids []ID) error {
	spotifyURL := c.baseURL + "me/shows?ids=" + strings.Join(toStringSlice(ids), ",")
	req, err := http.NewRequestWithContext(ctx, http.MethodPut, spotifyURL, nil)
	if err != nil {
		return err
	}

	return c.execute(req, nil, http.StatusOK)
}

// GetEpisode gets an episode from a show.
// API reference: https://developer.spotify.com/documentation/web-api/reference/get-an-episode
func (c *Client) GetEpisode(ctx context.Context, id string, opts ...RequestOption) (*EpisodePage, error) {
	spotifyURL := c.baseURL + "episodes/" + id
	if params := processOptions(opts...).urlParams.Encode(); params != "" {
		spotifyURL += "?" + params
	}

	var result EpisodePage

	err := c.get(ctx, spotifyURL, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}
