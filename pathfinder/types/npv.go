package types

const OpQueryNpvArtist Operation = "queryNpvArtist"

type QueryNpvArtistPayload struct {
	ArtistUri string `json:"artistUri"`
	TrackUri  string `json:"trackUri"`

	ContributorsLimit        *int  `json:"contributorsLimit,omitempty"`
	ContributorsOffset       *int  `json:"contributorsOffset,omitempty"`
	EnableRelatedVideos      *bool `json:"enableRelatedVideos,omitempty"`
	EnableRelatedAudioTracks *bool `json:"enableRelatedAudioTracks,omitempty"`
}

type QueryNpvArtistOption func(*QueryNpvArtistPayload)

func WithContributorsLimit(limit int) QueryNpvArtistOption {
	return func(p *QueryNpvArtistPayload) {
		p.ContributorsLimit = &limit
	}
}

func WithContributorsOffset(offset int) QueryNpvArtistOption {
	return func(p *QueryNpvArtistPayload) {
		p.ContributorsLimit = &offset
	}
}

func WithEnableRelatedVideos(enableRelatedVideos bool) QueryNpvArtistOption {
	return func(p *QueryNpvArtistPayload) {
		p.EnableRelatedVideos = &enableRelatedVideos
	}
}

func WithEnableRelatedAudioTracks(enableRelatedAudio bool) QueryNpvArtistOption {
	return func(p *QueryNpvArtistPayload) {
		p.EnableRelatedAudioTracks = &enableRelatedAudio
	}
}
