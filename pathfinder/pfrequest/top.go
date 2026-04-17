package pfrequest

import (
	"github.com/pyrorhythm/libspot/pathfinder/pfdomain"
)

type SearchTopRequest struct {
	SearchSuggestionsRequest

	SectionFilters []pfdomain.SectionFilter `json:"sectionFilters,omitempty"`
}

func (SearchTopRequest) Op() Operation {
	return OpSearchTop
}

func Top() *SearchTopRequest {
	return &SearchTopRequest{
		SearchSuggestionsRequest: *Suggestions(),
	}
}

func (st *SearchTopRequest) WithQuery(q string) *SearchTopRequest {
	st.Query = q
	return st
}

func (st *SearchTopRequest) WithCommons(o CommonsOpts) *SearchTopRequest {
	st.merge(o)
	return st
}

func (st *SearchTopRequest) WithSectionFilters(
	filters ...pfdomain.SectionFilter,
) *SearchTopRequest {
	st.SectionFilters = filters
	return st
}
