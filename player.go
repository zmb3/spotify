package spotify

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"strconv"
)

var (
	// ErrTemporarilyUnavailable When the device is temporarily unavailable the request will
	// return this error the client should retry the request after 5 seconds,
	// but no more than at most 5 retries.
	ErrTemporarilyUnavailable = errors.New("Device is temporarily unavailable")

	// ErrDeviceNotFound The requested device is not found
	ErrDeviceNotFound = errors.New("Device is not found")

	// ErrNotPremium The user making the request is non-premium
	ErrNotPremium = errors.New("User is not premium")
)

// PlayerDevice contains information about a device that a user can play music on
type PlayerDevice struct {
	// ID of the device. This may be empty.
	ID string `json:"id"`
	// IsActive If this device is the currently active device.
	IsActive bool `json:"is_active"`
	// IsRestricted Whether controlling this device is restricted. At present if
	// this is "true" then no Web API commands will be accepted by this device.
	IsRestricted bool `json:"is_restricted"`
	// Name The name of the device.
	Name string `json:"name"`
	// Type Device type, such as "Computer", "Smartphone" or "Speaker".
	Type string `json:"type"`
	// Volume The current volume in percent.
	Volume int `json:"volume_percent"`
}

// PlayerDeviceList is a list of devices
type PlayerDeviceList struct {
	PlayerDevices []PlayerDevice `json:"devices"`
}

// PlayerState contains information about the current playback.
type PlayerState struct {
	CurrentlyPlaying
	// Device The device that is currently active
	Device string `json:"device"`
	// ShuffleState Shuffle is on or off
	ShuffleState bool `json:"shuffle_state"`
	// RepeatState off, track, context
	RepeatState string `json:"repeat_state"`
}

// Context is the playback context
type Context struct {
	// ExternalURLs of the context, or null if not available.
	ExternalURLs map[string]string `json:"external_urls"`
	// Endpoint of the context, or null if not available.
	Endpoint string `json:"href"`
	// Type of the item's context. Can be one of album, artist or playlist.
	Type string `json:"type"`
	// URI of the context.
	URI URI `json:"uri"`
}

// CurrentlyPlaying contains the information about currently playing items
type CurrentlyPlaying struct {
	// Timestamp when data was fetched
	Timestamp int `json:"timestamp"`
	// Context current context
	Context Context `json:"context"`
	// Progress into the currently playing track.
	Progress string `json:"progress_ms"`
	// IsPlaying If something is currently playing.
	IsPlaying bool `json:"is_playing"`
	// The currently playing track. Can be null.
	Item *FullTrack `json:"Item"`
}

type Offset struct {
	//Position is zero based and can’t be negative.
	Position int `json:"position,omitempty"`
	// URI is a string representing the uri of the item to start at.
	URI URI `json:"uri,omitempty"`
}

type PlayOptions struct {
	// DeviceID The id of the device this command is targeting. If not
	// supplied, the user's currently active device is the target.
	DeviceID *string `json:"-"`
	// Context Spotify URI of the context to play.
	// Valid contexts are albums, artists & playlists.
	Context URI `json:"context,omitempty"`
	// URIs Array of the Spotify track URIs to play
	URIs []URI `json:"uris,omitempty"`
	// Offset Indicates from where in the context playback should start.
	// Only available when context corresponds to an album or playlist
	// object, or when the URIs parameter is used.
	Offset Offset `json:"offset,omitempty"`
}

func decodeStatusError(resp *http.Response) error {
	switch resp.StatusCode {
	case http.StatusNoContent:
		return nil
	case http.StatusAccepted:
		return ErrTemporarilyUnavailable
	case http.StatusNotFound:
		return ErrDeviceNotFound
	case http.StatusForbidden:
		return ErrNotPremium
	}
	return decodeError(resp.Body)
}

// PlayerDevices information about available devices for the current user.
// This call requires authorization.
//
// Requires the ScopeUserReadPlaybackState scope in order to read information
func (c *Client) PlayerDevices() (*PlayerDeviceList, error) {
	resp, err := c.http.Get(baseAddress + "me/player/devices")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, decodeError(resp.Body)
	}
	var result PlayerDeviceList
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// PlayerState gets information about the playing state for the current user
// This call requires authorization.
//
// Requires the ScopeUserReadPlaybackState scope in order to read information
func (c *Client) PlayerState() (*PlayerState, error) {
	return c.PlayerStateOpt(nil)
}

// PlayerStateOpt is like PlayerState, but it accepts additional
// options for sorting and filtering the results.
func (c *Client) PlayerStateOpt(opt *Options) (*PlayerState, error) {
	spotifyURL := baseAddress + "me/player"
	if opt != nil {
		v := url.Values{}
		if opt.Country != nil {
			v.Set("market", *opt.Country)
		}
		if params := v.Encode(); params != "" {
			spotifyURL += "?" + params
		}
	}
	resp, err := c.http.Get(spotifyURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, decodeError(resp.Body)
	}
	var result PlayerState
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// PlayerCurrentlyPlaying gets information about the currently playing status
// for the current user.  This call requires authorization.
//
// Requires the ScopeUserReadCurrentlyPlaying scope or the ScopeUserReadPlaybackState scope
// in order to read information
func (c *Client) PlayerCurrentlyPlaying() (*CurrentlyPlaying, error) {
	return c.PlayerCurrentlyPlayingOpt(nil)
}

// PlayerCurrentlyPlaying is like PlayerCurrentlyPlaying, but it accepts additional
// options for sorting and filtering the results.
func (c *Client) PlayerCurrentlyPlayingOpt(opt *Options) (*CurrentlyPlaying, error) {
	spotifyURL := baseAddress + "me/player/currently-playing"
	if opt != nil {
		v := url.Values{}
		if opt.Country != nil {
			v.Set("market", *opt.Country)
		}
		if params := v.Encode(); params != "" {
			spotifyURL += "?" + params
		}
	}
	resp, err := c.http.Get(spotifyURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, decodeError(resp.Body)
	}
	var result CurrentlyPlaying
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// TransferPlayback transfers playback to a new device and determine if
// it should start playing. This call requires authorization.
//
// Note that a value of false for the play parameter when also transferring
// to another device_id will not pause playback. To ensure that playback is
// paused on the new device you should send a pause command to the currently
// active device before transferring to the new device_id.
//
// Requires the ScopeUserModifyPlaybackState in order to modify the player state
func (c *Client) TransferPlayback(device_ids []string, play bool) error {
	reqData := struct {
		DeviceIDs []string `json:"device_ids"`
		Play      bool     `json:"play"`
	}{
		DeviceIDs: device_ids,
		Play:      play,
	}

	buf := new(bytes.Buffer)
	err := json.NewEncoder(buf).Encode(reqData)
	if err != nil {
		return err
	}
	req, err := http.NewRequest(http.MethodPut, baseAddress+"me/player", buf)
	if err != nil {
		return err
	}
	resp, err := c.http.Do(req)
	if err != nil {
		return err
	}
	return decodeStatusError(resp)
}

// Play Start a new context or resume current playback on the user's active
// device. This call requires authorization.
//
// Requires the ScopeUserModifyPlaybackState in order to modify the player state
func (c *Client) Play() error {
	return c.PlayOpt(nil)
}

// PlayOpt Like Play but with more options
func (c *Client) PlayOpt(opt *PlayOptions) error {
	spotifyURL := baseAddress + "me/player/play"
	buf := new(bytes.Buffer)

	if opt != nil {
		v := url.Values{}
		if opt.DeviceID != nil {
			v.Set("device_id", *opt.DeviceID)
		}
		if params := v.Encode(); params != "" {
			spotifyURL += "?" + params
		}

		err := json.NewEncoder(buf).Encode(opt)
		if err != nil {
			return err
		}
	}
	req, err := http.NewRequest(http.MethodPut, spotifyURL, buf)
	if err != nil {
		return err
	}
	resp, err := c.http.Do(req)
	if err != nil {
		return err
	}
	return decodeStatusError(resp)
}

// Pause Playback on the user's currently active device.
// This call requires authorization.
//
// Requires the ScopeUserModifyPlaybackState in order to modify the player state
func (c *Client) Pause() error {
	return c.PauseOpt(nil)
}

// PauseOpt Like Pause but with more options
//
// Only expects PlayOptions.DeviceID, all other options will be ignored
func (c *Client) PauseOpt(opt *PlayOptions) error {
	spotifyURL := baseAddress + "me/player/pause"

	if opt != nil {
		v := url.Values{}
		if opt.DeviceID != nil {
			v.Set("device_id", *opt.DeviceID)
		}
		if params := v.Encode(); params != "" {
			spotifyURL += "?" + params
		}
	}
	req, err := http.NewRequest(http.MethodPut, spotifyURL, nil)
	if err != nil {
		return err
	}
	resp, err := c.http.Do(req)
	if err != nil {
		return err
	}
	return decodeStatusError(resp)
}

// Next skips to the next track in the user's queue in the user's
// currently active device. This call requires authorization.
//
// Requires the ScopeUserModifyPlaybackState in order to modify the player state
func (c *Client) Next() error {
	return c.NextOpt(nil)
}

// NextOpt Like Next but with more options
//
// Only expects PlayOptions.DeviceID, all other options will be ignored
func (c *Client) NextOpt(opt *PlayOptions) error {
	spotifyURL := baseAddress + "me/player/next"

	if opt != nil {
		v := url.Values{}
		if opt.DeviceID != nil {
			v.Set("device_id", *opt.DeviceID)
		}
		if params := v.Encode(); params != "" {
			spotifyURL += "?" + params
		}
	}
	req, err := http.NewRequest(http.MethodPost, spotifyURL, nil)
	if err != nil {
		return err
	}
	resp, err := c.http.Do(req)
	if err != nil {
		return err
	}
	return decodeStatusError(resp)
}

// Previous skips to the Previous track in the user's queue in the user's
// currently active device. This call requires authorization.
//
// Requires the ScopeUserModifyPlaybackState in order to modify the player state
func (c *Client) Previous() error {
	return c.NextOpt(nil)
}

// PreviousOpt Like Previous but with more options
//
// Only expects PlayOptions.DeviceID, all other options will be ignored
func (c *Client) PreviousOpt(opt *PlayOptions) error {
	spotifyURL := baseAddress + "me/player/previous"

	if opt != nil {
		v := url.Values{}
		if opt.DeviceID != nil {
			v.Set("device_id", *opt.DeviceID)
		}
		if params := v.Encode(); params != "" {
			spotifyURL += "?" + params
		}
	}
	req, err := http.NewRequest(http.MethodPost, spotifyURL, nil)
	if err != nil {
		return err
	}
	resp, err := c.http.Do(req)
	if err != nil {
		return err
	}
	return decodeStatusError(resp)
}

// Seek to the given position in the user’s currently playing track.
// This call requires authorization.
//
// The position in milliseconds to seek to. Must be a positive number.
// Passing in a position that is greater than the length of the track
// will cause the player to start playing the next song.
//
// Requires the ScopeUserModifyPlaybackState in order to modify the player state
func (c *Client) Seek(position int) error {
	return c.SeekOpt(position, nil)
}

// SeekOpt Like Seek but with more options
//
// Only expects PlayOptions.DeviceID, all other options will be ignored
func (c *Client) SeekOpt(position int, opt *PlayOptions) error {
	return c.playerFuncWithOpt(
		"me/player/seek",
		url.Values{
			"position_ms": []string{strconv.FormatInt(int64(position), 10)},
		},
		opt,
	)
}

// Repeat Set the repeat mode for the user's playback.
// This call requires authorization.
//
// Options are repeat-track, repeat-context, and off.
//
// Requires the ScopeUserModifyPlaybackState in order to modify the player state
func (c *Client) Repeat(state string) error {
	return c.RepeatOpt(state, nil)
}

// RepeatOpt Like Repeat but with more options
//
// Only expects PlayOptions.DeviceID, all other options will be ignored
func (c *Client) RepeatOpt(state string, opt *PlayOptions) error {
	return c.playerFuncWithOpt(
		"me/player/repeat",
		url.Values{
			"state": []string{state},
		},
		opt,
	)
}

// Volume set the volume for the user's current playback device.
// This call requires authorization.
//
// Percent is must be a value from 0 to 100 inclusive.
//
// Requires the ScopeUserModifyPlaybackState in order to modify the player state
func (c *Client) Volume(percent int) error {
	return c.VolumeOpt(percent, nil)
}

// VolumeOpt Like Volume but with more options
//
// Only expects PlayOptions.DeviceID, all other options will be ignored
func (c *Client) VolumeOpt(percent int, opt *PlayOptions) error {
	return c.playerFuncWithOpt(
		"me/player/volume",
		url.Values{
			"volume_percent": []string{strconv.FormatInt(int64(percent), 10)},
		},
		opt,
	)
}

// Shuffle switches shuffle on or off for user's playback.
//
// Requires the ScopeUserModifyPlaybackState in order to modify the player state
func (c *Client) Shuffle(shuffle bool) error {
	return c.ShuffleOpt(shuffle, nil)
}

// ShuffleOpt Like Shuffle but with more options
//
// Only expects PlayOptions.DeviceID, all other options will be ignored
func (c *Client) ShuffleOpt(shuffle bool, opt *PlayOptions) error {
	return c.playerFuncWithOpt(
		"me/player/shuffle",
		url.Values{
			"state": []string{strconv.FormatBool(shuffle)},
		},
		opt,
	)
}

func (c *Client) playerFuncWithOpt(urlSuffix string, values url.Values, opt *PlayOptions) error {
	spotifyURL := baseAddress + urlSuffix

	if opt != nil {
		if opt.DeviceID != nil {
			values.Set("device_id", *opt.DeviceID)
		}
	}

	if params := values.Encode(); params != "" {
		spotifyURL += "?" + params
	}

	req, err := http.NewRequest(http.MethodPut, spotifyURL, nil)
	if err != nil {
		return err
	}
	resp, err := c.http.Do(req)
	if err != nil {
		return err
	}
	return decodeStatusError(resp)
}
