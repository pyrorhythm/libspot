package libspot

import (
	"encoding/hex"
	"fmt"
	"math/big"
	"regexp"
	"strings"
)

var uriRegexp = regexp.MustCompile("^spotify:([a-z]+):([0-9a-zA-Z]{21,22})$")

type SpotifyIdType string

const (
	SpotifyIdTypeTrack    SpotifyIdType = "track"
	SpotifyIdTypeEpisode  SpotifyIdType = "episode"
	SpotifyIdTypePlaylist SpotifyIdType = "playlist"
)

// SpotifyId is a Spotify object identifier (the 16-byte "gid") together with
// its type. It converts freely between the hex, base62 and URI forms.
type SpotifyId struct {
	typ SpotifyIdType
	id  []byte
}

func (id SpotifyId) Type() SpotifyIdType { return id.typ }
func (id SpotifyId) Id() []byte          { return id.id }
func (id SpotifyId) Hex() string         { return hex.EncodeToString(id.id) }
func (id SpotifyId) Base62() string      { return GidToBase62(id.id) }
func (id SpotifyId) Uri() string         { return fmt.Sprintf("spotify:%s:%s", id.typ, id.Base62()) }
func (id SpotifyId) String() string      { return id.Uri() }

// GidToBase62 encodes a 16-byte gid as its zero-padded 22-char base62 form.
func GidToBase62(id []byte) string {
	s := new(big.Int).SetBytes(id).Text(62)
	return strings.Repeat("0", 22-len(s)) + s
}

func SpotifyIdFromGid(typ SpotifyIdType, id []byte) SpotifyId {
	return SpotifyId{typ, id}
}

func SpotifyIdFromBase62(typ SpotifyIdType, id string) (*SpotifyId, error) {
	var i big.Int
	if _, ok := i.SetString(id, 62); !ok {
		return nil, fmt.Errorf("failed decoding base62: %s", id)
	}

	return &SpotifyId{typ, i.FillBytes(make([]byte, 16))}, nil
}

func SpotifyIdFromUri(uri string) (*SpotifyId, error) {
	matches := uriRegexp.FindStringSubmatch(uri)
	if len(matches) == 0 {
		return nil, fmt.Errorf("invalid uri: %s", uri)
	}

	return SpotifyIdFromBase62(SpotifyIdType(matches[1]), matches[2])
}
