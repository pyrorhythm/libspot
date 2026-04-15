package domaintypes

type AlbumSnippet struct {
	ID             string           `json:"id"`
	URI            string           `json:"uri"`
	Name           string           `json:"name"`
	CoverArt       *Image           `json:"coverArt"`
	VisualIdentity *VisualIdentity  `json:"visualIdentity"`
}

type VisualIdentityImage struct {
	ExtractedColorSet *ExtractedColorSet `json:"extractedColorSet"`
}

type AlbumResponseWrapper struct {
	URI            string                `json:"uri"`
	Name           string                `json:"name"`
	Type           AlbumResponseType     `json:"type"`
	Artists        ItemList[ArtistSnippet] `json:"artists"`
	CoverArt       *Image                `json:"coverArt"`
	Date           *DateSnippet          `json:"date"`
	Playability    *Playability          `json:"playability"`
	VisualIdentity *VisualIdentity       `json:"visualIdentity"`
}
