package spotify

import (
	"net/http"
	"reflect"
	"testing"
)

const fieldsDifferTemplate = "Actual response is not the same as expected response on field %s"

var expected = AudioAnalysis{
	Bars: []Marker{
		{
			Start:      251.98282,
			Duration:   0.29765,
			Confidence: 0.652,
		},
	},
	Beats: []Marker{
		{
			Start:      251.98282,
			Duration:   0.29765,
			Confidence: 0.652,
		},
	},
	Meta: AnalysisMeta{
		AnalyzerVersion: "4.0.0",
		Platform:        "Linux",
		DetailedStatus:  "OK",
		StatusCode:      0,
		Timestamp:       1456010389,
		AnalysisTime:    9.1394,
		InputProcess:    "libvorbisfile L+R 44100->22050",
	},
	Sections: []Section{
		{
			Marker: Marker{
				Start:      237.02356,
				Duration:   18.32542,
				Confidence: 1,
			},
			Loudness:                -20.074,
			Tempo:                   98.253,
			TempoConfidence:         0.767,
			Key:                     5,
			KeyConfidence:           0.327,
			Mode:                    1,
			ModeConfidence:          0.566,
			TimeSignature:           4,
			TimeSignatureConfidence: 1,
		},
	},
	Segments: []Segment{
		{
			Marker: Marker{
				Start:      252.15601,
				Duration:   3.19297,
				Confidence: 0.522,
			},
			LoudnessStart:   -23.356,
			LoudnessMaxTime: 0.06971,
			LoudnessMax:     -18.121,
			LoudnessEnd:     -60,
			Pitches:         []float64{0.709, 0.092, 0.196, 0.084, 0.352, 0.134, 0.161, 1, 0.17, 0.161, 0.211, 0.15},
			Timbre:          []float64{23.312, -7.374, -45.719, 294.874, 51.869, -79.384, -89.048, 143.322, -4.676, -51.303, -33.274, -19.037},
		},
	},
	Tatums: []Marker{
		{
			Start:      251.98282,
			Duration:   0.29765,
			Confidence: 0.652,
		},
	},
	Track: AnalysisTrack{
		NumSamples:              100,
		Duration:                255.34898,
		SampleMD5:               "",
		OffsetSeconds:           0,
		WindowSeconds:           0,
		AnalysisSampleRate:      22050,
		AnalysisChannels:        1,
		EndOfFadeIn:             0,
		StartOfFadeOut:          251.73333,
		Loudness:                -11.84,
		Tempo:                   98.002,
		TempoConfidence:         0.423,
		TimeSignature:           4,
		TimeSignatureConfidence: 1,
		Key:              5,
		KeyConfidence:    0.36,
		Mode:             0,
		ModeConfidence:   0.414,
		CodeString:       "eJxVnAmS5DgOBL-ST-B9_P9j4x7M6qoxW9tpsZQSCeI...",
		CodeVersion:      3.15,
		EchoprintString:  "eJzlvQmSHDmStHslxw4cB-v9j_A-tahhVKV0IH9...",
		EchoprintVersion: 4.12,
		SynchString:      "eJx1mIlx7ToORFNRCCK455_YoE9Dtt-vmrKsK3EBsTY...",
		SynchVersion:     1,
		RhythmString:     "eJyNXAmOLT2r28pZQuZh_xv7g21Iqu_3pCd160xV...",
		RhythmVersion:    1,
	},
}

func TestAudioAnalysis(t *testing.T) {
	c, s := testClientFile(http.StatusOK, "test_data/get_audio_analysis.txt")
	defer s.Close()

	analysis, err := c.GetAudioAnalysis("foo")
	if err != nil {
		t.Error(err)
	}

	if !reflect.DeepEqual(analysis.Bars, expected.Bars) {
		t.Errorf(fieldsDifferTemplate, "Bars")
	}

	if !reflect.DeepEqual(analysis.Beats, expected.Beats) {
		t.Errorf(fieldsDifferTemplate, "Beats")
	}

	if !reflect.DeepEqual(analysis.Meta, expected.Meta) {
		t.Errorf(fieldsDifferTemplate, "Meta")
	}

	if !reflect.DeepEqual(analysis.Sections, expected.Sections) {
		t.Errorf(fieldsDifferTemplate, "Sections")
	}

	if !reflect.DeepEqual(analysis.Segments, expected.Segments) {
		t.Errorf(fieldsDifferTemplate, "Segments")
	}

	if !reflect.DeepEqual(analysis.Track, expected.Track) {
		t.Errorf(fieldsDifferTemplate, "Track")
	}

	if !reflect.DeepEqual(analysis.Tatums, expected.Tatums) {
		t.Errorf(fieldsDifferTemplate, "Tatums")
	}
}
