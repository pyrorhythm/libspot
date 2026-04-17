package pfdomain

type Podcast struct {
	URI       string     `json:"uri"`
	Name      string     `json:"name"`
	Publisher *Publisher `json:"publisher"`
	MediaType MediaType  `json:"mediaType"`
	CoverArt  *ImageRaw  `json:"coverArt"`
}
