package metadata

type OriginalAudioItem struct {
	OriginalAudio OriginalAudio `json:"original_audio"`
}

type Track struct {
	Gid                            string                `json:"gid"`
	Name                           string                `json:"name"`
	Album                          Album                 `json:"album"`
	Artist                         []ArtistShort         `json:"artist"`
	Number                         int                   `json:"number"`
	DiscNumber                     int                   `json:"disc_number"`
	Duration                       int                   `json:"duration"`
	Popularity                     int                   `json:"popularity"`
	ExternalId                     []ExternalId          `json:"external_id"`
	EarliestLiveTimestamp          int                   `json:"earliest_live_timestamp"`
	HasLyrics                      bool                  `json:"has_lyrics"`
	Licensor                       Uuid                  `json:"licensor"`
	LanguageOfPerformance          []string              `json:"language_of_performance"`
	OriginalAudio                  OriginalAudio         `json:"original_audio"`
	OriginalTitle                  string                `json:"original_title"`
	ArtistWithRole                 []ArtistWithRole      `json:"artist_with_role"`
	CanonicalUri                   string                `json:"canonical_uri"`
	ContentAuthorizationAttributes string                `json:"content_authorization_attributes"`
	AudioFormats                   []OriginalAudioItem   `json:"audio_formats"`
	MediaType                      string                `json:"media_type"`
	ImplementationDetails          ImplementationDetails `json:"implementation_details"`
}

func (t Track) Type() MdType {
	return TypeTrack
}
