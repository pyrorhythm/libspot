package pfdomain

type PlaylistPreviewItems struct {
	Uri          string          `json:"_uri"`
	PreviewItems ItemList[Track] `json:"previewItems"`
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
