package spotify

import (
	"net/url"
	"strconv"
)

type Option func(*optionSet)

type optionSet struct {
	urlParams url.Values
}

func (oS optionSet) URLParams() url.Values {
	return oS.urlParams
}

// Limit sets the number of entries that a request should return
func Limit(amount int) Option {
	return func(oS *optionSet) {
		oS.urlParams.Set("limit", strconv.Itoa(amount))
	}
}

// Country enables track re-linking
func Country(code string) Option {
	return func(oS *optionSet) {
		oS.urlParams.Set("market", code)
	}
}

// Offset sets the index of the first entry to return
func Offset(amount int) Option {
	return func(oS *optionSet) {
		oS.urlParams.Set("offset", strconv.Itoa(amount))
	}
}

// Timerange sets t
func Timerange(timerange string) Option {
	return func(oS *optionSet) {
		oS.urlParams.Set("time_range", timerange+"_term")
	}
}

func processOptions(options ...Option) optionSet {
	oS := optionSet{}
	for _, option := range options {
		option(&oS)
	}

	return oS
}
