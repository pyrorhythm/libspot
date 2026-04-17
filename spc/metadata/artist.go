package metadata

type TopTrack struct {
	Country string `json:"country"`
	Track   []Gid  `json:"track"`
}

type ReleaseItem struct {
	Album []Gid `json:"album"`
}

type Portrait struct {
	Image []Image `json:"image"`
}

type Biography struct {
	Text          string     `json:"text"`
	PortraitGroup []Portrait `json:"portrait_group"`
}

type StartYear struct {
	StartYear int `json:"start_year"`
}

type Artist struct {
	ArtistShort

	Popularity     int           `json:"popularity"`
	TopTrack       []TopTrack    `json:"top_track"`
	AlbumGroup     []ReleaseItem `json:"album_group"`
	SingleGroup    []ReleaseItem `json:"single_group"`
	AppearsOnGroup []ReleaseItem `json:"appears_on_group"`
	Biography      []Biography   `json:"biography"`
	ActivityPeriod []StartYear   `json:"activity_period"`
	PortraitGroup  []Portrait    `json:"portrait_group"`
}

func (Artist) Type() MdType {
	return TypeArtist
}

type ArtistWithRole struct {
	ArtistGid  string `json:"artist_gid"`
	ArtistName string `json:"artist_name"`
	Role       string `json:"role"`
}
