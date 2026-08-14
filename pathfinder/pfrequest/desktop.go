package pfrequest

// SearchDesktopRequest is the aggregated search the desktop client uses: one
// query returns every bucket, including podcasts, whose dedicated operation
// hash the server no longer accepts.
type SearchDesktopRequest struct {
	SearchCommonsRequest

	SearchTerm string `json:"searchTerm"`
}

func (SearchDesktopRequest) Op() Operation {
	return OpSearchDesktop
}

func SearchDesktop(term string) *SearchDesktopRequest {
	return &SearchDesktopRequest{
		SearchCommonsRequest: defaultSearchCommons(),
		SearchTerm:           term,
	}
}

func SearchDesktopFromOpts(opts *BadgeRequestOpts) *SearchDesktopRequest {
	return &SearchDesktopRequest{
		SearchCommonsRequest: opts.SearchCommonsRequest,
		SearchTerm:           opts.SearchTerm,
	}
}

func (s *SearchDesktopRequest) WithCommons(o CommonsOpts) *SearchDesktopRequest {
	s.merge(o)
	return s
}
