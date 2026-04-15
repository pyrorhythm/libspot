package pfrequest

import "fmt"

type BadgeSearchPayload struct {
	*SearchPayloadCommons

	Kind       BadgeOperation `json:"-"`
	SearchTerm string         `json:"searchTerm"`
}

func (b BadgeSearchPayload) OperationType() (Operation, error) {
	if b.Kind.Valid() {
		return Operation(b.Kind), nil
	}

	return "", fmt.Errorf("invalid badge search operation kind: %v", b.Kind)
}

type BadgeSearchOption func(*BadgeSearchPayload)

func WithKind(kind BadgeOperation) BadgeSearchOption {
	return func(p *BadgeSearchPayload) {
		p.Kind = kind
	}
}
