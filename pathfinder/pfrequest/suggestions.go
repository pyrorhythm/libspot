package pfrequest

type SearchSuggestionsRequest struct {
	SearchCommonsRequest

	Query string `json:"query"`
}

func (SearchSuggestionsRequest) Op() Operation {
	return OpSearchSuggestions
}

func Suggestions() *SearchSuggestionsRequest {
	return &SearchSuggestionsRequest{SearchCommonsRequest: defaultSearchCommons()}
}

func (s *SearchSuggestionsRequest) WithQuery(q string) *SearchSuggestionsRequest {
	s.Query = q
	return s
}

func (s *SearchSuggestionsRequest) WithCommons(o CommonsOpts) *SearchSuggestionsRequest {
	s.merge(o)
	return s
}
