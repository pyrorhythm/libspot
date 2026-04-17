package pfrequest

type SearchSuggestionsRequest struct {
	SearchCommonsRequest

	Query string `json:"query"`
}

func (SearchSuggestionsRequest) Op() Operation {
	return OpSearchSuggestions
}

func Suggestions(q string) *SearchSuggestionsRequest {
	return &SearchSuggestionsRequest{
		SearchCommonsRequest: defaultSearchCommons(),
		Query:                q,
	}
}

func (s *SearchSuggestionsRequest) WithCommons(o CommonsOpts) *SearchSuggestionsRequest {
	s.merge(o)
	return s
}
