package pfrequest

// HomeRequest TODO
type HomeRequest struct {
	HomeEndUserIntegration string `json:"homeEndUserIntegration"`
	TimeZone               string `json:"timeZone"`
	SpT                    string `json:"sp_t"`
	Facet                  string `json:"facet"`
	SectionItemsLimit      int    `json:"sectionItemsLimit"`
}

func (HomeRequest) Op() Operation {
	return OpHome
}
