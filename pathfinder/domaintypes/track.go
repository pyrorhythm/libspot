package domaintypes

type TrackResponseWrapper struct {
	ID             string               `json:"id"`
	TrackMediaType TrackMediaType       `json:"trackMediaType"`
	Name           string               `json:"name"`
	Album          *AlbumSnippet        `json:"albumOfTrack"`
	Artists        ItemList[ArtistSnippet] `json:"artists"`
	Associations   *AssociationsV3      `json:"associationsV3"`
	ContentRating  *ContentRating       `json:"contentRating"`
	Duration       *Duration            `json:"duration"`
	Playability    *Playability         `json:"playability"`
	URI            string               `json:"uri"`
	VisualIdentity *VisualIdentity      `json:"visualIdentity"`
}
