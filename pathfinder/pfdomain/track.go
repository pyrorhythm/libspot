package pfdomain

type Track struct {
	Id             string                  `json:"id"`
	Uri            string                  `json:"uri"`
	TrackMediaType TrackMediaType          `json:"trackMediaType"`
	Name           string                  `json:"name"`
	Album          *AlbumSnippet           `json:"albumOfTrack"`
	Artists        ItemList[ArtistSnippet] `json:"artists"`
	Associations   *AssociationsV3         `json:"associationsV3"`
	ContentRating  *ContentRating          `json:"contentRating"`
	Duration       *Duration               `json:"duration"`
	Playability    *Playability            `json:"playability"`
	Previews       *Previews               `json:"previews"`
	VisualIdentity *VisualIdentity         `json:"visualIdentity"`
}

type Previews struct {
	ItemList[Url] `json:"audioPreviews"`
}

type TrackFromAlbum struct {
	TrackFromAlbumPayload

	Uid string `json:"uid"`
}

type TrackFromAlbumPayload struct {
	URI            string                  `json:"uri"`
	Name           string                  `json:"name"`
	Duration       Duration                `json:"duration"`
	Artists        ItemList[ArtistSnippet] `json:"artists"`
	ContentRating  ContentRating           `json:"contentRating"`
	Playability    Playability             `json:"playability"`
	AssociationsV3 AssociationsV3          `json:"associationsV3"`

	DiscNumber  int `json:"discNumber"`
	TrackNumber int `json:"trackNumber"`

	Playcount string `json:"playcount"`
	Saved     bool   `json:"saved"`

	RelinkingInformation any `json:"relinkingInformation"`
}
