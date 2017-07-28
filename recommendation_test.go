package spotify

import (
	"net/url"
	"testing"
)

func TestGetRecommendations(t *testing.T) {
	// test data corresponding to Spotify Console web API sample
	client, server := testClientFile(200, "test_data/recommendations.txt")
	defer server.Close()

	seeds := Seeds{
		Artists: []ID{"4NHQUGzhtTLFvgF5SZesLK"},
		Tracks:  []ID{"0c6xIDDpzE81m2q797ordA"},
		Genres:  []string{"classical", "country"},
	}
	country := "ES"
	limit := 10
	opts := Options{
		Country: &country,
		Limit:   &limit,
	}
	recommendations, err := client.GetRecommendations(seeds, nil, &opts)
	if err != nil {
		t.Fatal(err)
	}
	if len(recommendations.Tracks) != 10 {
		t.Error("Expected 10 recommended tracks")
	}
	if recommendations.Tracks[0].Artists[0].Name != "Heinrich Isaac" {
		t.Error("Expected the artist of the first recommended track to be Heinrich Isaac")
	}
}

func TestSetSeedValues(t *testing.T) {
	expectedValues := "seed_artists=4NHQUGzhtTLFvgF5SZesLK%2C5PHQUGzhtTUIvgF5SZesGY&seed_genres=classical%2Ccountry"
	v := url.Values{}
	seeds := Seeds{
		Artists: []ID{"4NHQUGzhtTLFvgF5SZesLK", "5PHQUGzhtTUIvgF5SZesGY"},
		Genres:  []string{"classical", "country"},
	}
	setSeedValues(seeds, v)
	actualValues := v.Encode()
	if actualValues != expectedValues {
		t.Errorf("Expected seed values to be %s but got %s", expectedValues, actualValues)
	}
}

func TestSetTrackAttributesValues(t *testing.T) {
	expectedValues := "max_duration_ms=200&min_duration_ms=20&min_energy=0.45&target_acousticness=0.27&target_duration_ms=160"
	v := url.Values{}
	ta := NewTrackAttributes().
		MaxDuration(200).
		MinDuration(20).
		TargetDuration(160).
		MinEnergy(0.45).
		TargetAcousticness(0.27)

	setTrackAttributesValues(ta, v)
	actualValues := v.Encode()
	if actualValues != expectedValues {
		t.Errorf("Expected track attributes values to be %s but got %s", expectedValues, actualValues)
	}
}

func TestSetEmptyTrackAttributesValues(t *testing.T) {
	expectedValues := ""
	v := url.Values{}
	setTrackAttributesValues(nil, v)
	actualValues := v.Encode()
	if actualValues != expectedValues {
		t.Errorf("Expected track attributes values to be empty but got %s", actualValues)
	}
}
