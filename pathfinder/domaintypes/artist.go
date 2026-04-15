package domaintypes

type ArtistSnippet struct {
	URI     string                `json:"uri"`
	Profile *ArtistSnippetProfile `json:"profile"`
}

type ItemList[T any] struct {
	Items []*T `json:"items"`
}

type ArtistSnippetProfile struct {
	Name string `json:"name"`
}

type ArtistResponseWrapper struct {
	ArtistSnippet

	VisualIdentity *VisualIdentity `json:"visualIdentity"`
	Visuals        *ArtistVisuals  `json:"visuals"`
}

type ArtistVisuals struct {
	AvatarImage *Image `json:"avatarImage"`
}
