package pfrequest

type FeedBaselineLookupRequest struct {
	Uris []string `json:"uris"`
}

func (FeedBaselineLookupRequest) Op() Operation {
	return OpFeedBaselineLookup
}

type ImageUris struct {
	ImageUris []string `json:"imageUris"`
}

type FetchExtractedColorsRequest ImageUris

func (FetchExtractedColorsRequest) Op() Operation {
	return OpFetchExtractedColors
}

type DynamicColorsRequest ImageUris

func (DynamicColorsRequest) Op() Operation {
	return OpGetDynamicColorsByUris
}
