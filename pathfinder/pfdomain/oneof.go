package pfdomain

import (
	"encoding/json"
	"log/slog"

	"github.com/pkg/errors"
	"github.com/pyrorhythm/fn/bjs"
	"github.com/valyala/fastjson"
)

var (
	ErrWrongType       = errors.New("responsetypes: wrong type — use Is*() before accessor")
	ErrUnknownTypename = errors.New("responsetypes: unknown __typename")
)

type caseOneof string

func (c caseOneof) String() string {
	return string(c)
}

const (
	caseCompletion caseOneof = "SearchAutoCompleteEntity"
	caseArtist     caseOneof = "ArtistResponseWrapper"
	caseTrack      caseOneof = "TrackResponseWrapper"
	caseAlbum      caseOneof = "AlbumResponseWrapper"
	caseEpisode    caseOneof = "EpisodeResponseWrapper"
	casePodcast    caseOneof = "PodcastResponseWrapper"
	casePlaylist   caseOneof = "PlaylistResponseWrapper"
	caseUser       caseOneof = "UserResponseWrapper"
	caseGenre      caseOneof = "GenreResponseWrapper"
)

type OneofMatched struct {
	Item *Oneof `json:"item"`

	MatchedFields []any `json:"matchedFields"`
}

type Oneof struct {
	typname caseOneof

	completion *SearchCompletion
	artist     *Artist
	track      *Track
	album      *Album
	episode    *Episode
	podcast    *Podcast
	playlist   *Playlist
	user       *User
	genre      *Genre
}

func (o *Oneof) HasArtist() bool     { return o.typname == caseArtist }
func (o *Oneof) HasTrack() bool      { return o.typname == caseTrack }
func (o *Oneof) HasAlbum() bool      { return o.typname == caseAlbum }
func (o *Oneof) HasEpisode() bool    { return o.typname == caseEpisode }
func (o *Oneof) HasPodcast() bool    { return o.typname == casePodcast }
func (o *Oneof) HasPlaylist() bool   { return o.typname == casePlaylist }
func (o *Oneof) HasUser() bool       { return o.typname == caseUser }
func (o *Oneof) HasGenre() bool      { return o.typname == caseGenre }
func (o *Oneof) HasCompletion() bool { return o.typname == caseCompletion }

func (o *Oneof) GetArtist() *Artist               { return o.artist }
func (o *Oneof) GetTrack() *Track                 { return o.track }
func (o *Oneof) GetAlbum() *Album                 { return o.album }
func (o *Oneof) GetEpisode() *Episode             { return o.episode }
func (o *Oneof) GetPodcast() *Podcast             { return o.podcast }
func (o *Oneof) GetPlaylist() *Playlist           { return o.playlist }
func (o *Oneof) GetUser() *User                   { return o.user }
func (o *Oneof) GetGenre() *Genre                 { return o.genre }
func (o *Oneof) GetCompletion() *SearchCompletion { return o.completion }

func (o *Oneof) UnmarshalJSON(data []byte) error {
	val, err := fastjson.ParseBytes(data)
	if err != nil {
		return errors.Wrap(err, "failed to parse json obj")
	}

	typname := string(val.Get("__typename").GetStringBytes())
	// slog.Debug("TopResultHitOneof.typename", "typname", typname)

	o.typname = caseOneof(typname)

	if !val.Exists("data") {
		return nil
	}
	payload := val.Get("data").MarshalTo(nil)

	switch o.typname {
	case caseCompletion:
		o.completion, err = bjs.Unmarshal[SearchCompletion](payload)
	case caseArtist:
		o.artist, err = bjs.Unmarshal[Artist](payload)
	case caseTrack:
		o.track, err = bjs.Unmarshal[Track](payload)
	case caseAlbum:
		o.album, err = bjs.Unmarshal[Album](payload)
	case caseEpisode:
		o.episode, err = bjs.Unmarshal[Episode](payload)
	case casePodcast:
		o.podcast, err = bjs.Unmarshal[Podcast](payload)
	case casePlaylist:
		o.playlist, err = bjs.Unmarshal[Playlist](payload)
	case caseGenre:
		o.genre, err = bjs.Unmarshal[Genre](payload)
	case caseUser:
		o.user, err = bjs.Unmarshal[User](payload)
	default:
		return errors.Wrapf(ErrUnknownTypename, "typename=%q", o.typname)
	}

	if err != nil {
		return errors.Wrapf(err, "failed to unmarshal type %s", o.typname)
	}

	slog.Debug("topResultHit", "type", o.typname)
	return nil
}

func (o *Oneof) MarshalJSON() ([]byte, error) {
	toMarshal := map[string]any{
		"__typename": o.typname.String(),
	}

	switch o.typname {
	case caseCompletion:
		toMarshal["data"] = o.completion
	case caseArtist:
		toMarshal["data"] = o.artist
	case caseTrack:
		toMarshal["data"] = o.track
	case caseAlbum:
		toMarshal["data"] = o.album
	case caseEpisode:
		toMarshal["data"] = o.episode
	case casePodcast:
		toMarshal["data"] = o.podcast
	case casePlaylist:
		toMarshal["data"] = o.playlist
	case caseGenre:
		toMarshal["data"] = o.genre
	case caseUser:
		toMarshal["data"] = o.user
	default:
		toMarshal["data"] = nil
	}

	return json.Marshal(toMarshal)
}
