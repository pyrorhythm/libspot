package domaintypes

type ImageSource struct {
	Height int    `json:"height"`
	Url    string `json:"url"`
	Width  int    `json:"width"`
}

type Image struct {
	ExtractedColors ExtractedColors `json:"extractedColors"`
	Sources         []*ImageSource  `json:"sources"`
}

type ExtractedColors struct {
	ColorDark ColorDark `json:"colorDark"`
}
