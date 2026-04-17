package pfdomain

import (
	"time"
)

type AssociationsV3 struct {
	AudioAssociations TotalCount `json:"audioAssociations"`
	VideoAssociations TotalCount `json:"videoAssociations"`
}

type Copyright struct {
	Text string `json:"text"`
	Type string `json:"type"`
}

type SharingInfo struct {
	ShareID  *string `json:"shareId"`
	ShareURL string  `json:"shareUrl"`
}

type Disc struct {
	Number int        `json:"number"`
	Tracks TotalCount `json:"tracks"`
}

type ContentRating struct {
	Label ContentRatingEnum `json:"label"`
}

type Uri struct {
	Uri string `json:"uri"`
}

type Url struct {
	Url string `json:"url"`
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

type PlayedState struct {
	PlayPositionMilliseconds int    `json:"playPositionMilliseconds"`
	State                    string `json:"state"`
}

type Publisher struct {
	Name string `json:"name"`
}

type Date struct {
	IsoString time.Time `json:"isoString"`
	Precision string    `json:"precision"`
}

type Restrictions struct {
	PaywallContent bool `json:"paywallContent"`
}
