package spotify

// SimpleArtist contains basic info about an artist.
type SimpleArtist struct {
	// The name of the artist.
	Name string `json:"name"`
	// The Spotify ID for the artist.
	ID ID `json:"id"`
	// The Spotify URI for the artist.
	URI URI `json:"uri"`
	// A link to the Web API enpoint providing
	// full details of the artist.
	Endpoint string `json:"href"`
	// Known external URLs for this artist.
	ExternalURLs ExternalURL `json:"external_urls"`
}

// FullArtist provides extra artist data in addition
// to what is provided by SimpleArtist.
type FullArtist struct {
	SimpleArtist
	// The popularity of the artist.  The value will be
	// between 0 and 100, with 100 being the most popular.
	// The artist's popularity is calculated from the
	// popularity of all of the artist's tracks.
	Popularity int `json:"popularity"`
	// A list of genres the artist is associated with.
	// For example, "Prog Rock" or "Post-Grunge".  If
	// not yet classified, the slice is empty.
	Genres []string `json:"genres"`
	// Information about followers of the artist.
	Followers Followers
	// Images of the artist in various sizes, widest first.
	Images []Image `json:"images"`
}

func (a *SimpleArtist) String() string {
	return "SimpleArtist: " + a.Name
}

func (a *FullArtist) String() string {
	return "FullArtist: " + a.Name
}
