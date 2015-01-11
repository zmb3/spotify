package spotify

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"
)

// SimpleTrack contains basic info about a track.
type SimpleTrack struct {
	// The artists who performed the track.
	Artists []SimpleArtist `json:"artists"`
	// A list of the countries in which the track
	// can be played, identified by their ISO 3166-1
	// alpha-2 codes.
	AvailableMarkets []string `json:"available_markets"`
	// The disc number (usually 1 unless the album
	// consists of more than one disc).
	DiscNumber int `json:"disc_number"`
	// The length of the track, in milliseconds. TODO: time package?
	Duration int `json:"duration_ms"`
	// Whether or not the track has explicit lyrics.
	// true => yes, it does; false => no, it does not.
	Explicit bool `json:"explicit"`
	// External URLs for this track.
	ExternalURLs ExternalURL `json:"external_urls"`
	// A link to the Web API endpoint providing full
	// details for this track.
	Endpoint string `json:"href"`
	// The Spotify ID for the track.
	ID ID `json:"id"`
	// The name of the track
	Name string `json:"name"`
	// A URL to a 30 second preview (MP3) of the track.
	PreviewURL string `json:"preview_url"`
	// The number of the track.  If an album has several
	// discs, the track number is the number on the specified
	// DiscNumber.
	TrackNumber int `json:"track_number"`
	// The Spotify URI for the track.
	URI URI `json:"uri"`
}

// FullTrack provides extra track data in addition
// to what is provided by SimpleTrack.
type FullTrack struct {
	SimpleTrack
	// Known external IDs for the track.
	ExternalIDs ExternalID `json:"exernal_ids"`
	// Popularity of the trakc.  The value will be
	// between 0 and 100, with 100 being the most
	// popular.  The popularity is calculated from
	// both total plays and most recent plays.
	Popularity int `json:"popularity"`
}

// PlaylistTrack contains info about a track
// in a playlist.
type PlaylistTrack struct {
	// The date and time the track was added to the playlist.
	// TODO: very old playlists may return null here.
	AddedAt Timestamp `json:"added_at"`
	// The Spotify user who added the track to the playlist.
	// TODO: very old playlists may return null here
	AddedBy User `json:"added_by"`
	// Information about the track.
	Track FullTrack `json:"track"`
}

// SavedTrack provides info about a track saved
// to a user's account.
type SavedTrack struct {
}

// FindTrack gets spotify catalog information for
// a single track identified by its unique Spotify ID.
func (c *Client) FindTrack(id ID) (*FullTrack, error) {
	uri := baseAddress + "tracks/" + string(id)
	resp, err := c.http.Get(uri)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		var e struct {
			E Error `json:"error"`
		}
		err = json.NewDecoder(resp.Body).Decode(&e)
		if err != nil {
			return nil, errors.New("spotify: Couldn't decode error object")
		}
		return nil, e.E
	}
	var t FullTrack
	err = json.NewDecoder(resp.Body).Decode(&t)
	if err != nil {
		return nil, err
	}
	return &t, nil
}

// FindTracks gets Spotify catalog information for multiple
// tracks based on their Spotify IDs.  It supports up to 50
// tracks in a single call.  Tracks are returned in the order
// requested.  If a track is not found, that position in the
// result will be nil.  Duplicate ids in the query will result
// in duplicate tracks in the result.
func (c *Client) FindTracks(ids ...ID) ([]*FullTrack, error) {
	if len(ids) > 50 {
		return nil, errors.New("spotify: FindTracks supports up to 50 tracks")
	}
	uri := baseAddress + "tracks?ids=" + strings.Join(toStringSlice(ids), ",")
	resp, err := c.http.Get(uri)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		var e struct {
			E Error `json:"error"`
		}
		err = json.NewDecoder(resp.Body).Decode(&e)
		if err != nil {
			return nil, errors.New("spotify: couldn't decode error")
		}
		return nil, e.E
	}

	var t struct {
		Tracks []*FullTrack `jsosn:"tracks"`
	}
	err = json.NewDecoder(resp.Body).Decode(&t)
	if err != nil {
		return nil, errors.New("spotify:  couldn't decode tracks")
	}
	return t.Tracks, nil
}
