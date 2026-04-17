package pfdomain

type Episode struct {
	URI                   string          `json:"uri"`
	Name                  string          `json:"name"`
	ContentRating         *ContentRating  `json:"contentRating"`
	CoverArt              *Image          `json:"coverArt"`
	Description           string          `json:"description"`
	Duration              *Duration       `json:"duration"`
	GatedEntityRelations  []any           `json:"gatedEntityRelations"`
	MediaTypes            []MediaType     `json:"mediaTypes"`
	Playability           *Playability    `json:"playability"`
	PlayedState           *PlayedState    `json:"playedState"`
	PodcastV2             Data[Podcast]   `json:"podcastV2"`
	ReleaseDate           *Date           `json:"releaseDate"`
	Restrictions          *Restrictions   `json:"restrictions"`
	VideoPreviewThumbnail any             `json:"videoPreviewThumbnail"`
	VisualIdentity        *VisualIdentity `json:"visualIdentity"`
}
