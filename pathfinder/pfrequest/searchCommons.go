package pfrequest

type SearchCommonsRequest struct {
	NumberOfTopResults             *int  `json:"numberOfTopResults,omitempty"`
	Offset                         *int  `json:"offset,omitempty"`
	Limit                          *int  `json:"limit,omitempty"`
	IncludePreReleases             *bool `json:"includePreReleases,omitempty"`
	IncludeArtistHasConcertsField  *bool `json:"includeArtistHasConcertsField,omitempty"`
	IncludeAudiobooks              *bool `json:"includeAudiobooks,omitempty"`
	IncludeAuthors                 *bool `json:"includeAuthors,omitempty"`
	IncludeEpisodeContentRatingsV2 *bool `json:"includeEpisodeContentRatingsV2,omitempty"`
}

func defaultSearchCommons() SearchCommonsRequest {
	return SearchCommonsRequest{
		NumberOfTopResults:             new(30),
		Offset:                         new(0),
		Limit:                          new(30),
		IncludePreReleases:             new(true),
		IncludeArtistHasConcertsField:  new(false),
		IncludeAudiobooks:              new(true),
		IncludeAuthors:                 new(true),
		IncludeEpisodeContentRatingsV2: new(false),
	}
}

type CommonsOpts SearchCommonsRequest

func (sc *SearchCommonsRequest) merge(o CommonsOpts) {
	if o.NumberOfTopResults != nil {
		sc.NumberOfTopResults = o.NumberOfTopResults
	}
	if o.Offset != nil {
		sc.Offset = o.Offset
	}
	if o.Limit != nil {
		sc.Limit = o.Limit
	}
	if o.IncludePreReleases != nil {
		sc.IncludePreReleases = o.IncludePreReleases
	}
	if o.IncludeArtistHasConcertsField != nil {
		sc.IncludeArtistHasConcertsField = o.IncludeArtistHasConcertsField
	}
	if o.IncludeAudiobooks != nil {
		sc.IncludeAudiobooks = o.IncludeAudiobooks
	}
	if o.IncludeAuthors != nil {
		sc.IncludeAuthors = o.IncludeAuthors
	}
	if o.IncludeEpisodeContentRatingsV2 != nil {
		sc.IncludeEpisodeContentRatingsV2 = o.IncludeEpisodeContentRatingsV2
	}
}

func Commons() CommonsOpts { return CommonsOpts{} }

func (o CommonsOpts) WithLimit(n int) CommonsOpts        { o.Limit = &n; return o }
func (o CommonsOpts) WithOffset(n int) CommonsOpts       { o.Offset = &n; return o }
func (o CommonsOpts) WithTopResults(n int) CommonsOpts   { o.NumberOfTopResults = &n; return o }
func (o CommonsOpts) WithAudiobooks(v bool) CommonsOpts  { o.IncludeAudiobooks = &v; return o }
func (o CommonsOpts) WithAuthors(v bool) CommonsOpts     { o.IncludeAuthors = &v; return o }
func (o CommonsOpts) WithPrereleases(v bool) CommonsOpts { o.IncludePreReleases = &v; return o }
func (o CommonsOpts) WithArtistConcerts(v bool) CommonsOpts {
	o.IncludeArtistHasConcertsField = &v
	return o
}

func (o CommonsOpts) WithEpisodeRatings(v bool) CommonsOpts {
	o.IncludeEpisodeContentRatingsV2 = &v
	return o
}
