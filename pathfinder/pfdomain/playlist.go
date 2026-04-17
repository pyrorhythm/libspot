package pfdomain

type PlaylistPreviewItems struct {
	Uri          string          `json:"_uri"`
	PreviewItems ItemList[Track] `json:"previewItems"`
}

type T struct {
	Items []struct {
		Typename string `json:"__typename"`
		Data     struct {
			Typename     string `json:"__typename"`
			AlbumOfTrack struct {
				CoverArt struct {
					ExtractedColors struct {
						ColorDark struct {
							Hex string `json:"hex"`
						} `json:"colorDark"`
					} `json:"extractedColors"`
					Sources []struct {
						Height int    `json:"height"`
						Url    string `json:"url"`
						Width  int    `json:"width"`
					} `json:"sources"`
				} `json:"coverArt"`
			} `json:"albumOfTrack"`
			Canvas interface{} `json:"canvas"`
			Name   string      `json:"name"`

			Uri string `json:"uri"`
		} `json:"data"`
	} `json:"items"`
}

type Playlist struct {
	URI            string          `json:"uri"`
	Name           string          `json:"name"`
	OwnerV2        Data[User]      `json:"ownerV2"`
	Description    string          `json:"description"`
	Attributes     []any           `json:"attributes"`
	Format         string          `json:"format"`
	Images         ItemList[Image] `json:"images"`
	VisualIdentity *VisualIdentity `json:"visualIdentity"`
}
