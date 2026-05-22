package connect

import (
	"encoding/json"

	"github.com/pkg/errors"
	"github.com/valyala/fastjson"
	"pyrorhythm.dev/fn/bjs"
)

var (
	ErrWrongMediaType = errors.New("connect: wrong media type — use Has*() before accessor")
	ErrUnknownMedia   = errors.New("connect: unknown media type")
)

type mediaKind string

const (
	mediaTrack   mediaKind = "track"
	mediaEpisode mediaKind = "episode"
)

// Track is a catalog track object (full or partial from connect-state).
type Track struct {
	ID      string   `json:"id,omitempty"`
	URI     string   `json:"uri"`
	Name    string   `json:"name"`
	Album   string   `json:"album,omitempty"`
	Artists []string `json:"artists,omitempty"`
}

// Episode is a podcast episode object (full or partial from connect-state).
type Episode struct {
	ID   string `json:"id,omitempty"`
	URI  string `json:"uri"`
	Name string `json:"name"`
}

// MediaOneof discriminates track vs episode payloads by uri kind or JSON type.
type MediaOneof struct {
	kind mediaKind

	track   *Track
	episode *Episode
}

func (o *MediaOneof) HasTrack() bool   { return o != nil && o.kind == mediaTrack }
func (o *MediaOneof) HasEpisode() bool { return o != nil && o.kind == mediaEpisode }

func (o *MediaOneof) GetTrack() *Track     { return o.track }
func (o *MediaOneof) GetEpisode() *Episode { return o.episode }
func (o *MediaOneof) Kind() string         { return string(o.kind) }

// PartialTrack wraps scalar connect-state fields as a track MediaOneof.
func PartialTrack(id, uri, name, album string, artists []string) *MediaOneof {
	return &MediaOneof{
		kind: mediaTrack,
		track: &Track{
			ID: id, URI: uri, Name: name, Album: album, Artists: artists,
		},
	}
}

// PartialEpisode wraps scalar connect-state fields as an episode MediaOneof.
func PartialEpisode(id, uri, name string) *MediaOneof {
	return &MediaOneof{
		kind:    mediaEpisode,
		episode: &Episode{ID: id, URI: uri, Name: name},
	}
}

func (o *MediaOneof) Item() MediaItem {
	if o == nil {
		return MediaItem{}
	}
	switch o.kind {
	case mediaTrack:
		if t := o.track; t != nil {
			return MediaItem{
				ID: t.ID, URI: t.URI, Name: t.Name, Type: string(mediaTrack),
				Album: t.Album, Artists: t.Artists,
			}
		}
	case mediaEpisode:
		if e := o.episode; e != nil {
			return MediaItem{
				ID: e.ID, URI: e.URI, Name: e.Name, Type: string(mediaEpisode),
			}
		}
	}
	return MediaItem{}
}

func (o *MediaOneof) UnmarshalJSON(data []byte) error {
	val, err := fastjson.ParseBytes(data)
	if err != nil {
		return errors.Wrap(err, "connect: parse media json")
	}
	typ := string(val.Get("type").GetStringBytes())
	if typ == "" {
		typ = typeFromURI(string(val.Get("uri").GetStringBytes()))
	}
	o.kind = mediaKind(typ)
	switch o.kind {
	case mediaTrack:
		o.track, err = bjs.Unmarshal[Track](data)
	case mediaEpisode:
		o.episode, err = bjs.Unmarshal[Episode](data)
	default:
		return errors.Wrapf(ErrUnknownMedia, "type=%q", typ)
	}
	return errors.Wrapf(err, "connect: unmarshal %s", o.kind)
}

func (o *MediaOneof) MarshalJSON() ([]byte, error) {
	if o == nil {
		return []byte("null"), nil
	}
	return json.Marshal(o.Item())
}
