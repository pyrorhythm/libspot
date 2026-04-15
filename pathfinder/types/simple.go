package types

type SuggestionsPayload struct {
	*SearchPayloadCommons

	Query string `json:"query"`
}

type FeedBaselineLookupPayload struct {
	Uris []string `json:"uris"`
}

type FetchExtractedColorsPayload struct {
	ImageUris []string `json:"imageUris"`
}
