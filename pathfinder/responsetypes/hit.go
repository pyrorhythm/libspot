package responsetypes

import (
	"log/slog"

	"github.com/pkg/errors"
	"github.com/pyrorhythm/fn/bjs"
	"github.com/pyrorhythm/libspot/pathfinder/domaintypes"
	"github.com/valyala/fastjson"
)

var (
	ErrWrongType       = errors.New("responsetypes: wrong type — use Is*() before accessor")
	ErrUnknownTypename = errors.New("responsetypes: unknown __typename")
)

type topResultHitOneof string

const (
	oneofTopResultHitCaseAutoComplete topResultHitOneof = "SearchAutoCompleteEntity"
	oneofTopResultHitCaseArtist       topResultHitOneof = "ArtistResponseWrapper"
	oneofTopResultHitCaseTrack        topResultHitOneof = "TrackResponseWrapper"
	oneofTopResultHitCaseAlbum        topResultHitOneof = "AlbumResponseWrapper"
)

type TopResultHit struct {
	Item          *TopResultHitOneof `json:"item"`
	MatchedFields []any              `json:"matchedFields"`
}

// TopResultHitOneof holds the raw item JSON and lazily unmarshals it into a typed entity.
type TopResultHitOneof struct {
	typname topResultHitOneof

	searchAutoComplete *domaintypes.SearchAutoCompleteEntity
	artist             *domaintypes.ArtistResponseWrapper
	track              *domaintypes.TrackResponseWrapper
	album              *domaintypes.AlbumResponseWrapper
}

func (h *TopResultHitOneof) HasSearchAutoComplete() bool {
	return h.typname == oneofTopResultHitCaseAutoComplete
}

func (h *TopResultHitOneof) HasArtist() bool { /*  */
	return h.typname == oneofTopResultHitCaseArtist
}

func (h *TopResultHitOneof) HasTrack() bool {
	return h.typname == oneofTopResultHitCaseTrack
}

func (h *TopResultHitOneof) HasAlbum() bool {
	return h.typname == oneofTopResultHitCaseAlbum
}

func (h *TopResultHitOneof) GetSearchAutoComplete() *domaintypes.SearchAutoCompleteEntity {
	return h.searchAutoComplete
}

func (h *TopResultHitOneof) GetArtist() *domaintypes.ArtistResponseWrapper {
	return h.artist
}

func (h *TopResultHitOneof) GetTrack() *domaintypes.TrackResponseWrapper {
	return h.track
}

func (h *TopResultHitOneof) GetAlbum() *domaintypes.AlbumResponseWrapper {
	return h.album
}

func (h *TopResultHitOneof) UnmarshalJSON(data []byte) error {
	val, err := fastjson.ParseBytes(data)
	if err != nil {
		return errors.Wrap(err, "failed to parse json obj")
	}

	typname := string(val.Get("__typename").GetStringBytes())
	// slog.Debug("TopResultHitOneof.typename", "typname", typname)

	h.typname = topResultHitOneof(typname)
	payload := val.Get("data").MarshalTo(nil)

	switch h.typname {
	case oneofTopResultHitCaseAutoComplete:
		h.searchAutoComplete, err = bjs.Unmarshal[domaintypes.SearchAutoCompleteEntity](payload)
	case oneofTopResultHitCaseArtist:
		h.artist, err = bjs.Unmarshal[domaintypes.ArtistResponseWrapper](payload)
	case oneofTopResultHitCaseTrack:
		h.track, err = bjs.Unmarshal[domaintypes.TrackResponseWrapper](payload)
	case oneofTopResultHitCaseAlbum:
		h.album, err = bjs.Unmarshal[domaintypes.AlbumResponseWrapper](payload)
	default:
		return errors.Wrapf(ErrUnknownTypename, "typename=%q", h.typname)
	}

	if err != nil {
		return errors.Wrapf(err, "failed to unmarshal type %s", h.typname)
	}

	slog.Debug("topResultHit", "type", h.typname)
	return nil
}
