package pfdomain

import (
	"encoding/json"

	"github.com/pkg/errors"
	"github.com/valyala/fastjson"
)

type PlayabilityReason string

const (
	PlayabilityReasonPlayable PlayabilityReason = "PLAYABLE"
)

type TrackMediaType string

const (
	TrackMediaTypeAudio TrackMediaType = "AUDIO"
)

type AlbumResponseType string

const AlbumResponseTypeAlbum AlbumResponseType = "ALBUM"

type Chip string

func (c *Chip) UnmarshalJSON(bytes []byte) error {
	v, err := fastjson.ParseBytes(bytes)
	if err != nil {
		return errors.Wrap(err, "failed to parse chip")
	}

	chipVal := v.Get("typeName").MarshalTo(nil)
	switch Chip(chipVal) {
	case ChipAlbums, ChipArtists,
		ChipAudiobooks, ChipAuthors,
		ChipEpisodes, ChipGenres,
		ChipPlaylists, ChipPodcasts,
		ChipTracks, ChipUsers:
		*c = Chip(chipVal)
		return nil
	default:
		return errors.Errorf("invalid chip type: %s", chipVal)
	}
}

type Chips []Chip

func (c *Chips) UnmarshalJSON(bytes []byte) error {
	v, err := fastjson.ParseBytes(bytes)
	if err != nil {
		return errors.Wrap(err, "failed to parse chip")
	}

	chipsVal := v.Get("items").MarshalTo(nil)
	var arr []Chip
	err = json.Unmarshal(chipsVal, &arr)
	if err != nil {
		return errors.Wrap(err, "failed to unmarshal chips")
	}
	*c = arr
	return nil
}

const (
	ChipTracks     Chip = "TRACKS"
	ChipAlbums     Chip = "ALBUMS"
	ChipArtists    Chip = "ARTISTS"
	ChipPlaylists  Chip = "PLAYLISTS"
	ChipEpisodes   Chip = "EPISODES"
	ChipPodcasts   Chip = "PODCASTS"
	ChipUsers      Chip = "USERS"
	ChipAudiobooks Chip = "AUDIOBOOKS"
	ChipAuthors    Chip = "AUTHORS"
	ChipGenres     Chip = "GENRES"
)
