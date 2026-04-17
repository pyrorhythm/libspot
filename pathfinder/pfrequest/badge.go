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

func (b *BadgeSearchRequest) WithTerm(t string) *BadgeSearchRequest {
	b.SearchTerm = t
	return b
}

func (b *BadgeSearchRequest) WithCommons(o CommonsOpts) *BadgeSearchRequest {
	b.SearchCommonsRequest.merge(o)
	return b
}
