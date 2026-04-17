package pfdomain

type ImageSource struct {
	Height int    `json:"height"`
	Url    string `json:"url"`
	Width  int    `json:"width"`
}

type ImageRaw struct {
	Sources []*ImageSource `json:"sources"`
}

type ImageRawWithFormat struct {
	ImageRaw

	Format ImageFormat `json:"imageFormat"`
}

type Image struct {
	ImageRaw

	ExtractedColors *ExtractedColor `json:"extractedColors"`
}

type ThumbnailImage struct {
	ImageID     string               `json:"imageId"`
	ImageIDType ImageV2Type          `json:"imageIdType"`
	Sources     []ImageRawWithFormat `json:"sources"`
}

type AvatarImage struct {
	AvatarImage *Image `json:"avatarImage"`
}

type VisualIdentity struct {
	SixteenByNineCoverImage *VisualIdentityImage `json:"sixteenByNineCoverImage"`
	SquareCoverImage        *VisualIdentityImage `json:"squareCoverImage"`
}

type VisualIdentityImage struct {
	ExtractedColorSet *ExtractedColorSet `json:"extractedColorSet"`
}
