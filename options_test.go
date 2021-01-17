package spotify

import (
	"testing"
)

func TestOptions(t *testing.T) {
	t.Parallel()

	resultSet := processOptions(
		Offset(1),
		Limit(13),
		Country(CountryUnitedKingdom),
		Timerange("long"),
	)

	expected := "limit=13&market=GB&offset=1&time_range=long_term"
	actual := resultSet.urlParams.Encode()
	if actual != expected {
		t.Errorf("Expected '%v', got '%v'", expected, actual)
	}
}
