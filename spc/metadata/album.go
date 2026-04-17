package metadata

type Disc struct {
	Number int   `json:"number"`
	Track  []Gid `json:"track"`
}

type Album struct {
	AlbumShort

	// Type                           string           `json:"type"`
	Popularity                     int              `json:"popularity"`
	ExternalId                     []ExternalId     `json:"external_id"`
	Disc                           []Disc           `json:"disc"`
	Copyright                      []Copyright      `json:"copyright"`
	OriginalTitle                  string           `json:"original_title"`
	EarliestLiveTimestamp          int              `json:"earliest_live_timestamp"`
	CanonicalUri                   string           `json:"canonical_uri"`
	ArtistWithRole                 []ArtistWithRole `json:"artist_with_role"`
	ContentAuthorizationAttributes string           `json:"content_authorization_attributes"`
}

func (a Album) Type() MdType {
	return TypeAlbum
}

type AlbumShort struct {
	Gid        string        `json:"gid"`
	Name       string        `json:"name"`
	Artist     []ArtistShort `json:"artist"`
	Label      string        `json:"label"`
	Date       Date          `json:"date"`
	CoverGroup CoverGroup    `json:"cover_group"`
	Licensor   Uuid          `json:"licensor"`
}
