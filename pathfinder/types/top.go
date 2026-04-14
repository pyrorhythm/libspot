package types

const OpTopSearch Operation = "searchTopResultsList"

type SectionFilter string

const (
	SFGeneric      SectionFilter = "GENERIC"
	SFVideoContent SectionFilter = "VIDEO_CONTENT"
)

type TopSearchPayload struct {
	*SearchPayloadCommons

	Query          string          `json:"query"`
	SectionFilters []SectionFilter `json:"sectionFilters,omitempty"`
	// TODO: provide an enum
}

type TopSearchOption func(payload *TopSearchPayload)

func WithSectionFilters(filters ...SectionFilter) TopSearchOption {
	return func(payload *TopSearchPayload) {
		payload.SectionFilters = filters
	}
}
