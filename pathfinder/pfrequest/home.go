package pfrequest

// HomePayload TODO
type HomePayload struct {
	HomeEndUserIntegration string `json:"homeEndUserIntegration"`
	TimeZone               string `json:"timeZone"`
	SpT                    string `json:"sp_t"`
	Facet                  string `json:"facet"`
	SectionItemsLimit      int    `json:"sectionItemsLimit"`
}
