package pfrequest

type SearchPayloadCommons struct {
	NumberOfTopResults             *int  `json:"numberOfTopResults,omitempty"`
	Offset                         *int  `json:"offset,omitempty"`
	Limit                          *int  `json:"limit,omitempty"`
	IncludePreReleases             *bool `json:"includePreReleases,omitempty"`
	IncludeArtistHasConcertsField  *bool `json:"includeArtistHasConcertsField,omitempty"`
	IncludeAudiobooks              *bool `json:"includeAudiobooks,omitempty"`
	IncludeAuthors                 *bool `json:"includeAuthors,omitempty"`
	IncludeEpisodeContentRatingsV2 *bool `json:"includeEpisodeContentRatingsV2,omitempty"`
}

type SearchCommonsOption func(*SearchPayloadCommons)

func WithLimit(n int) SearchCommonsOption {
	return func(r *SearchPayloadCommons) { r.Limit = &n }
}

func WithOffset(n int) SearchCommonsOption {
	return func(r *SearchPayloadCommons) { r.Offset = &n }
}

func WithTopResults(n int) SearchCommonsOption {
	return func(r *SearchPayloadCommons) { r.NumberOfTopResults = &n }
}

func WithAudiobooks(v bool) SearchCommonsOption {
	return func(r *SearchPayloadCommons) { r.IncludeAudiobooks = &v }
}

func WithArtistHasConcertsField(v bool) SearchCommonsOption {
	return func(r *SearchPayloadCommons) { r.IncludeArtistHasConcertsField = &v }
}

func WithAuthors(v bool) SearchCommonsOption {
	return func(r *SearchPayloadCommons) { r.IncludeAuthors = &v }
}

func WithPrereleases(v bool) SearchCommonsOption {
	return func(r *SearchPayloadCommons) { r.IncludePreReleases = &v }
}
