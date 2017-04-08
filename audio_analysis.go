package spotify

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// AudioAnalysis contains a detailed audio analysis for a single track identified by its unique Spotify ID.
// See https://developer.spotify.com/web-api/get-audio-analysis/
//
// Spotify's documentation is currently missing the object model for the AudioAnalysis.
// See https://github.com/spotify/web-api/issues/317
//
// Also see The Echo Nest documentation
// https://web.archive.org/web/20160528174915/http://developer.echonest.com/docs/v4/_static/AnalyzeDocumentation.pdf
type AudioAnalysis struct {
	Bars     []Measure          `json:"bars"`
	Beats    []Measure          `json:"beats"`
	Meta     AudioAnalysisMeta  `json:"meta"`
	Sections []Section          `json:"sections"`
	Segments []Segment          `json:"segments"`
	Tatums   []Measure          `json:"tatums"`
	Track    AudioAnalysisTrack `json:"track"`
}

// Measure represents beats, bars, tatums and are used in segments and sections descriptions.
type Measure struct {
	Start      float64 `json:"start"`
	Duration   float64 `json:"duration"`
	Confidence float64 `json:"confidence"`
}

type AudioAnalysisMeta struct {
	AnalyzerVersion string  `json:"analyzer_version"`
	Platform        string  `json:"platform"`
	DetailedStatus  string  `json:"detailed_status"`
	StatusCode      int     `json:"status"`
	Timestamp       int64   `json:"timestamp"`
	AnalysisTime    float64 `json:"analysis_time"`
	InputProcess    string  `json:"input_process"`
}

type Section struct {
	Measure
	Loudness                float64 `json:"loudness"`
	Tempo                   float64 `json:"tempo"`
	TempoConfidence         float64 `json:"tempo_confidence"`
	Key                     int     `json:"key"`
	KeyConfidence           float64 `json:"key_confidence"`
	Mode                    int     `json:"mode"`
	ModeConfidence          float64 `json:"mode_confidence"`
	TimeSignature           int     `json:"time_signature"`
	TimeSignatureConfidence float64 `json:"time_signature_confidence"`
}

type Segment struct {
	Measure
	LoudnessStart   float64   `json:"loudness_start"`
	LoudnessMaxTime float64   `json:"loudness_max_time"`
	LoudnessMax     float64   `json:"loudness_max"`
	LoudnessEnd     float64   `json:"loudness_end"`
	Pitches         []float64 `json:"pitches"`
	Timbre          []float64 `json:"timbre"`
}

type AudioAnalysisTrack struct {
	NumSamples              int64   `json:"num_samples"`
	Duration                float64 `json:"duration"`
	SampleMD5               string  `json:"sample_md5"`
	OffsetSeconds           int     `json:"offset_seconds"`
	WindowSeconds           int     `json:"window_seconds"`
	AnalysisSampleRate      int64   `json:"analysis_sample_rate"`
	AnalysisChannels        int     `json:"analysis_channels"`
	EndOfFadeIn             float64 `json:"end_of_fade_in"`
	StartOfFadeOut          float64 `json:"start_of_fade_out"`
	Loudness                float64 `json:"loudness"`
	Tempo                   float64 `json:"tempo"`
	TempoConfidence         float64 `json:"tempo_confidence"`
	TimeSignature           int     `json:"time_signature"`
	TimeSignatureConfidence float64 `json:"time_signature_confidence"`
	Key                     Key     `json:"key"`
	KeyConfidence           float64 `json:"key_confidence"`
	Mode                    Mode    `json:"mode"`
	ModeConfidence          float64 `json:"mode_confidence"`
	CodeString              string  `json:"codestring"`
	CodeVersion             float64 `json:"code_version"`
	EchoprintString         string  `json:"echoprintstring"`
	EchoprintVersion        float64 `json:"echoprint_version"`
	SynchString             string  `json:"synchstring"`
	SynchVersion            float64 `json:"synch_version"`
	RhythmString            string  `json:"rhythmstring"`
	RhythmVersion           float64 `json:"rhythm_version"`
}

// GetAudioAnalysis queries the Spotify web API for an audio analysis of a single track
// If an object is not found, a nil value is returned in the appropriate position.
// This call requires authorization.
func (c *Client) GetAudioAnalysis(id ID) (*AudioAnalysis, error) {
	url := fmt.Sprintf("%saudio-analysis/%s", baseAddress, id)

	resp, err := c.http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, decodeError(resp.Body)
	}

	temp := AudioAnalysis{}
	err = json.NewDecoder(resp.Body).Decode(&temp)
	if err != nil {
		return nil, err
	}

	return &temp, nil
}
