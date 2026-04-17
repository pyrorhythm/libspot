package pfrequest

type BadgeSearchRequest struct {
	SearchCommonsRequest

	Kind       BadgeOperation `json:"-"`
	SearchTerm string         `json:"searchTerm"`
}

func (b BadgeSearchRequest) Op() Operation {
	return Operation(b.Kind)
}

func BadgeSearch(kind BadgeOperation) *BadgeSearchRequest {
	return &BadgeSearchRequest{
		SearchCommonsRequest: defaultSearchCommons(),
		Kind:                 kind,
	}
}

func BadgeSearchFromOpts(kind BadgeOperation, opts *BadgeRequestOpts) *BadgeSearchRequest {
	return &BadgeSearchRequest{
		SearchCommonsRequest: opts.SearchCommonsRequest,
		Kind:                 kind,
	}
}

func (b *BadgeSearchRequest) WithTerm(t string) *BadgeSearchRequest {
	b.SearchTerm = t
	return b
}

func (b *BadgeSearchRequest) WithCommons(o CommonsOpts) *BadgeSearchRequest {
	b.SearchCommonsRequest.merge(o)
	return b
}

type BadgeRequestOpts struct {
    SearchCommonsRequest

	SearchTerm string
}

func BadgeOpts() *BadgeRequestOpts {
	return &BadgeRequestOpts{
		SearchCommonsRequest: defaultSearchCommons(),
	}
}

func (b *BadgeRequestOpts) WithTerm(t string) *BadgeRequestOpts {
	b.SearchTerm = t
	return b
}

func (b *BadgeRequestOpts) WithCommons(o CommonsOpts) *BadgeRequestOpts {
	b.SearchCommonsRequest.merge(o)
	return b
}