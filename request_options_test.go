package spotify

import (
	"testing"
)

func TestOptions(t *testing.T) {
	t.Parallel()

	resultSet := processOptions(
		After("example_id"),
		Country(CountryUnitedKingdom),
		Limit(13),
		Locale("en_GB"),
		Market(CountryArgentina),
		Offset(1),
		Timerange("long"),
		Timestamp("2000-11-02T13:37:00"),
	)

	expected := "after=example_id&country=GB&limit=13&locale=en_GB&market=AR&offset=1&time_range=long&timestamp=2000-11-02T13%3A37%3A00"
	actual := resultSet.urlParams.Encode()
	if actual != expected {
		t.Errorf("Expected '%v', got '%v'", expected, actual)
	}
}
