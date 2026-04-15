package pfresponse

import (
	"log/slog"

	"github.com/pkg/errors"
	"github.com/pyrorhythm/fn/bjs"
	"github.com/pyrorhythm/libspot/pathfinder/pfdomain"
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

	searchAutoComplete *pfdomain.SearchAutoCompleteEntity
	artist             *pfdomain.ArtistResponseWrapper
	track              *pfdomain.TrackResponseWrapper
	album              *pfdomain.AlbumResponseWrapper
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

func (h *TopResultHitOneof) GetSearchAutoComplete() *pfdomain.SearchAutoCompleteEntity {
	return h.searchAutoComplete
}

func (h *TopResultHitOneof) GetArtist() *pfdomain.ArtistResponseWrapper {
	return h.artist
}

func (h *TopResultHitOneof) GetTrack() *pfdomain.TrackResponseWrapper {
	return h.track
}

func (h *TopResultHitOneof) GetAlbum() *pfdomain.AlbumResponseWrapper {
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
		h.searchAutoComplete, err = bjs.Unmarshal[pfdomain.SearchAutoCompleteEntity](payload)
	case oneofTopResultHitCaseArtist:
		h.artist, err = bjs.Unmarshal[pfdomain.ArtistResponseWrapper](payload)
	case oneofTopResultHitCaseTrack:
		h.track, err = bjs.Unmarshal[pfdomain.TrackResponseWrapper](payload)
	case oneofTopResultHitCaseAlbum:
		h.album, err = bjs.Unmarshal[pfdomain.AlbumResponseWrapper](payload)
	default:
		return errors.Wrapf(ErrUnknownTypename, "typename=%q", h.typname)
	}

	if err != nil {
		return errors.Wrapf(err, "failed to unmarshal type %s", h.typname)
	}

	slog.Debug("topResultHit", "type", h.typname)
	return nil
}
