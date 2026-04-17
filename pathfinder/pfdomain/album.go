package pfdomain

type AlbumSnippet struct {
	ID             string          `json:"id"`
	URI            string          `json:"uri"`
	Name           string          `json:"name"`
	CoverArt       *Image          `json:"coverArt"`
	VisualIdentity *VisualIdentity `json:"visualIdentity"`
}

type Album struct {
	URI            string                  `json:"uri"`
	Name           string                  `json:"name"`
	Type           AlbumResponseType       `json:"type"`
	Artists        ItemList[ArtistSnippet] `json:"artists"`
	CoverArt       *Image                  `json:"coverArt"`
	Date           *DateSnippet            `json:"date"`
	Playability    *Playability            `json:"playability"`
	VisualIdentity *VisualIdentity         `json:"visualIdentity"`
}

type AlbumFromReleases struct {
	Name string `json:"name"`
	URI  string `json:"uri"`
}

type AlbumFromDiscography struct {
	ID          string      `json:"id"`
	URI         string      `json:"uri"`
	Name        string      `json:"name"`
	Type        string      `json:"type"`
	CoverArt    Image       `json:"coverArt"`
	Date        DateSnippet `json:"date"`
	Playability Playability `json:"playability"`
	SharingInfo SharingInfo `json:"sharingInfo"`
}

type (
	Discography struct {
		ItemList[AlbumFromDiscography] `json:"popularReleasesAlbums"`
	}

	MoreAlbumsByArtist struct {
		Discography `json:"discography"`
	}

	AlbumFull struct {
		URI                   string                           `json:"uri"`
		Name                  string                           `json:"name"`
		Type                  string                           `json:"type"`
		Date                  Date                             `json:"date"`
		Saved                 bool                             `json:"saved"`
		CoverArt              Image                            `json:"coverArt"`
		SharingInfo           SharingInfo                      `json:"sharingInfo"`
		VisualIdentity        VisualIdentity                   `json:"visualIdentity"`
		Playability           Playability                      `json:"playability"`
		WatchFeedEntrypoint   WatchFeedEntrypoint              `json:"watchFeedEntrypoint"`
		IsPreRelease          bool                             `json:"isPreRelease"`
		PreReleaseEndDateTime Date                             `json:"preReleaseEndDateTime"`
		CourtesyLine          string                           `json:"courtesyLine"`
		Label                 string                           `json:"label"`
		Copyright             ItemCountList[Copyright]         `json:"copyright"`
		Tracks                ItemCountList[TrackFromAlbum]    `json:"tracksV2"`
		Artists               ItemCountList[ArtistFromAlbum]   `json:"artists"`
		Discs                 ItemCountList[Disc]              `json:"discs"`
		Releases              ItemCountList[AlbumFromReleases] `json:"releases"`
		MoreAlbumsByArtist    ItemList[MoreAlbumsByArtist]     `json:"moreAlbumsByArtist"`
	}
)
