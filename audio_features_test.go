package spotify

import (
	"net/http"
	"testing"
)

var response = `
{
  "audio_features" : [ {
    "danceability" : 0.808,
    "energy" : 0.626,
    "key" : 7,
    "loudness" : -12.733,
    "mode" : 1,
    "speechiness" : 0.168,
    "acousticness" : 0.00187,
    "instrumentalness" : 0.159,
    "liveness" : 0.376,
    "valence" : 0.369,
    "tempo" : 123.990,
    "type" : "audio_features",
    "id" : "4JpKVNYnVcJ8tuMKjAj50A",
    "uri" : "spotify:track:4JpKVNYnVcJ8tuMKjAj50A",
    "track_href" : "https://api.spotify.com/v1/tracks/4JpKVNYnVcJ8tuMKjAj50A",
    "analysis_url" : "http://echonest-analysis.s3.amazonaws.com/TR/WhpYUARk1kNJ_qP0AdKGcDDFKOQTTgsOoINrqyPQjkUnbteuuBiyj_u94iFCSGzdxGiwqQ6d77f4QLL_8=/3/full.json?AWSAccessKeyId=AKIAJRDFEY23UEVW42BQ&Expires=1459290544&Signature=4P03WGLL1a/%2BXp90jcsLGMfFC3Y%3D",
    "duration_ms" : 535223,
    "time_signature" : 4
  }, {
    "danceability" : 0.457,
    "energy" : 0.815,
    "key" : 1,
    "loudness" : -7.199,
    "mode" : 1,
    "speechiness" : 0.0340,
    "acousticness" : 0.102,
    "instrumentalness" : 0.0319,
    "liveness" : 0.103,
    "valence" : 0.382,
    "tempo" : 96.083,
    "type" : "audio_features",
    "id" : "2NRANZE9UCmPAS5XVbXL40",
    "uri" : "spotify:track:2NRANZE9UCmPAS5XVbXL40",
    "track_href" : "https://api.spotify.com/v1/tracks/2NRANZE9UCmPAS5XVbXL40",
    "analysis_url" : "http://echonest-analysis.s3.amazonaws.com/TR/WhuQhwPDhmEg5TO4JjbJu0my-awIhk3eaXkRd1ofoJ7tXogPnMtbxkTyLOeHXu5Jke0FCIt52saKJyfPM=/3/full.json?AWSAccessKeyId=AKIAJRDFEY23UEVW42BQ&Expires=1459290544&Signature=Jsg/GexxC7v06Tq70coL/d2x7kI%3D",
    "duration_ms" : 187800,
    "time_signature" : 4
  }, null, {
    "danceability" : 0.281,
    "energy" : 0.402,
    "key" : 4,
    "loudness" : -17.921,
    "mode" : 1,
    "speechiness" : 0.0291,
    "acousticness" : 0.0734,
    "instrumentalness" : 0.830,
    "liveness" : 0.0593,
    "valence" : 0.0748,
    "tempo" : 115.700,
    "type" : "audio_features",
    "id" : "24JygzOLM0EmRQeGtFcIcG",
    "uri" : "spotify:track:24JygzOLM0EmRQeGtFcIcG",
    "track_href" : "https://api.spotify.com/v1/tracks/24JygzOLM0EmRQeGtFcIcG",
    "analysis_url" : "http://echonest-analysis.s3.amazonaws.com/TR/ehbkMg05Ck-FN7p3lV7vd8TUdBCvM6z5mgDiZRv6iSlw8P_b8GYBZ4PRAlOgTl3e5rS34_l3dZGDeYzH4=/3/full.json?AWSAccessKeyId=AKIAJRDFEY23UEVW42BQ&Expires=1459290544&Signature=09T3QyRucjrOMoMutRmdJKLJ7hI%3D",
    "duration_ms" : 497493,
    "time_signature" : 3
  } ]
}
`

func TestAudioFeatures(t *testing.T) {
	c, s := testClientString(http.StatusOK, response)
	defer s.Close()

	ids := []ID{
		"4JpKVNYnVcJ8tuMKjAj50A",
		"2NRANZE9UCmPAS5XVbXL40",
		"abc", // intentionally throw a bad one in
		"24JygzOLM0EmRQeGtFcIcG",
	}
	features, err := c.GetAudioFeatures()
	if err != nil {
		t.Error(err)
	}
	if len(features) != len(ids) {
		t.Errorf("Want %d results, got %d\n", len(ids), len(features))
	}
	if features[2] != nil {
		t.Errorf("Want nil result, got #%v\n", features[2])
	}
	if Key(features[0].Key) != G {
		t.Errorf("Want key G, got %v\n", features[0].Key)
	}
}
