package spotify

import (
	"net/url"
	"strconv"
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

// Country enables track re-linking
func Country(code string) RequestOption {
	return func(o *requestOptions) {
		o.urlParams.Set("market", code)
	}
}

// Offset sets the index of the first entry to return
func Offset(amount int) RequestOption {
	return func(o *requestOptions) {
		o.urlParams.Set("offset", strconv.Itoa(amount))
	}
}

type Range string

var (
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

func processOptions(options ...RequestOption) requestOptions {
	o := requestOptions{
		urlParams: url.Values{},
	}
	for _, opt := range options {
		opt(&o)
	}

	return o
}
