package pfrequest

type QueryNpvArtistPayload struct {
	ArtistUri string `json:"artistUri"`
	TrackUri  string `json:"trackUri"`

	ContributorsLimit        *int  `json:"contributorsLimit,omitempty"`
	ContributorsOffset       *int  `json:"contributorsOffset,omitempty"`
	EnableRelatedVideos      *bool `json:"enableRelatedVideos,omitempty"`
	EnableRelatedAudioTracks *bool `json:"enableRelatedAudioTracks,omitempty"`
}

func (QueryNpvArtistPayload) Op() Operation {
	return OpQueryNpvArtist
}

func QueryNpvArtist(artistUri, trackUri string) *QueryNpvArtistPayload {
	return &QueryNpvArtistPayload{ArtistUri: artistUri, TrackUri: trackUri}
}

func (o *QueryNpvArtistPayload) WithContributorsLimit(n int) *QueryNpvArtistPayload {
	o.ContributorsLimit = &n
	return o
}

func (o *QueryNpvArtistPayload) WithContributorsOffset(n int) *QueryNpvArtistPayload {
	o.ContributorsOffset = &n
	return o
}

func (o *QueryNpvArtistPayload) WithRelatedVideos(v bool) *QueryNpvArtistPayload {
	o.EnableRelatedVideos = &v
	return o
}

func (o *QueryNpvArtistPayload) WithRelatedAudioTracks(v bool) *QueryNpvArtistPayload {
	o.EnableRelatedAudioTracks = &v
	return o
}
