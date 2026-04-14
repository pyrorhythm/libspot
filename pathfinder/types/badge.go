package types

import "fmt"

type BadgeOperation Operation

const (
	BOTracks    BadgeOperation = "searchTracks"
	BOAlbums    BadgeOperation = "searchAlbums"
	BOArtists   BadgeOperation = "searchArtists"
	BOPlaylists BadgeOperation = "searchPlaylists"
	BOProfiles  BadgeOperation = "searchProfiles"
)

func (b BadgeOperation) String() string {
	return string(b)
}

func (b BadgeOperation) Valid() bool {
	switch b {
	case BOTracks, BOAlbums, BOArtists, BOPlaylists, BOProfiles:
		return true
	}

	return false
}

type BadgeSearchPayload struct {
	*SearchPayloadCommons

	Kind       BadgeOperation `json:"-"`
	SearchTerm string         `json:"searchTerm"`
}

func (b BadgeSearchPayload) OperationType() (Operation, error) {
	if !b.Kind.Valid() {
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
