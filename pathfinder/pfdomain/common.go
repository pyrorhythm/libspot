package pfdomain

type ContentRating struct {
	Label string `json:"label"`
}

type Duration struct {
	TotalMilliseconds int `json:"totalMilliseconds"`
}

type Playability struct {
	Playable bool              `json:"playable"`
	Reason   PlayabilityReason `json:"reason"`
}

type DateSnippet struct {
	Year int `json:"year"`
}
