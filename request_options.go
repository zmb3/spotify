package spotify

import (
	"net/url"
	"strconv"
	"strings"
)

type RequestOption func(*requestOptions)

type requestOptions struct {
	urlParams url.Values
}

// Limit sets the number of entries that a request should return
func Limit(amount int) RequestOption {
	return func(o *requestOptions) {
		o.urlParams.Set("limit", strconv.Itoa(amount))
	}
}

// Market enables track re-linking
func Market(code string) RequestOption {
	return func(o *requestOptions) {
		o.urlParams.Set("market", code)
	}
}

// Country enables a specific region to be specified for region-specific suggestions e.g popular playlists
// The Country option takes an ISO 3166-1 alpha-2 country code.  It can be
// used to ensure that the category exists for a particular country.
func Country(code string) RequestOption {
	return func(o *requestOptions) {
		o.urlParams.Set("country", code)
	}
}

// Locale enables a specific language to be used when returning results.
// The Locale argument is an ISO 639 language code and an ISO 3166-1 alpha-2
// country code, separated by an underscore.  It can be used to get the
// category strings in a particular language (for example: "es_MX" means
// get categories in Mexico, returned in Spanish).
func Locale(code string) RequestOption {
	return func(o *requestOptions) {
		o.urlParams.Set("locale", code)
	}
}

// Offset sets the index of the first entry to return
func Offset(amount int) RequestOption {
	return func(o *requestOptions) {
		o.urlParams.Set("offset", strconv.Itoa(amount))
	}
}

// Timestamp in ISO 8601 format (yyyy-MM-ddTHH:mm:ss).
// use this parameter to specify the user's local time to
// get results tailored for that specific date and time
// in the day.  If not provided, the response defaults to
// the current UTC time.
func Timestamp(ts string) RequestOption {
	return func(o *requestOptions) {
		o.urlParams.Set("timestamp", ts)
	}
}

// After is the last ID retrieved from the previous request. This allows pagination.
func After(after string) RequestOption {
	return func(o *requestOptions) {
		o.urlParams.Set("after", after)
	}
}

// Fields is a comma-separated list of the fields to return.
// See the JSON tags on the FullPlaylist struct for valid field options.
// For example, to get just the playlist's description and URI:
//    fields = "description,uri"
//
// A dot separator can be used to specify non-reoccurring fields, while
// parentheses can be used to specify reoccurring fields within objects.
// For example, to get just the added date and the user ID of the adder:
//    fields = "tracks.items(added_at,added_by.id)"
//
// Use multiple parentheses to drill down into nested objects, for example:
//    fields = "tracks.items(track(name,href,album(name,href)))"
//
// Fields can be excluded by prefixing them with an exclamation mark, for example;
//    fields = "tracks.items(track(name,href,album(!name,href)))"
func Fields(fields string) RequestOption {
	return func(o *requestOptions) {
		o.urlParams.Set("fields", fields)
	}
}

type Range string

const (
	// LongTermRange is calculated from several years of data, including new data where possible
	LongTermRange Range = "long_term"
	// MediumTermRange is approximately the last six months
	MediumTermRange Range = "medium_term"
	// ShortTermRange is approximately the last four weeks
	ShortTermRange Range = "short_term"
)

// Timerange sets the time period that spoty will use when returning information. Use LongTermRange, MediumTermRange
// and ShortTermRange to set the appropriate period.
func Timerange(timerange Range) RequestOption {
	return func(o *requestOptions) {
		o.urlParams.Set("time_range", string(timerange))
	}
}

type AdditionalType string

const (
	EpisodeAdditionalType = "episode"
	TrackAdditionalType   = "track"
)

// AdditionalTypes is a list of item types that your client supports besides the default track type.
// Valid types are: EpisodeAdditionalType and TrackAdditionalType.
func AdditionalTypes(types ...AdditionalType) RequestOption {
	strTypes := make([]string, len(types))
	for i, t := range types {
		strTypes[i] = string(t)
	}

	csv := strings.Join(strTypes, ",")

	return func(o *requestOptions) {
		o.urlParams.Set("additional_types", csv)
	}
}

func processOptions(options ...RequestOption) requestOptions {
	o := requestOptions{
		urlParams: url.Values{},
	}
	for _, opt := range options {
		opt(&o)
	}

	return o
}
