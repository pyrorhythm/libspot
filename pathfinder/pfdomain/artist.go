package pfdomain

type ArtistSnippet struct {
	URI     string                `json:"uri"`
	Profile *ArtistSnippetProfile `json:"profile"`
}

type ArtistSnippetProfile struct {
	Name string `json:"name"`
}

type Artist struct {
	ArtistSnippet

	VisualIdentity *VisualIdentity `json:"visualIdentity"`
	Visuals        *AvatarImage    `json:"visuals"`
}

type ArtistFromAlbum struct {
	Id string `json:"id"`

	ArtistSnippet //nolint

	SharingInfo *SharingInfo `json:"sharingInfo"`
	Visuals     *AvatarImage `json:"visuals"`
}
