package extendp

import (
	"github.com/pyrorhythm/libspot/spc/metadata"
)

const Path = "playlistextender/extendp"

type Request struct {
	PlaylistURI   string          `json:"playlistURI"`
	TrackIDs      []*metadata.Gid `json:"trackIDs"`
	ArtistIDs     []*metadata.Gid `json:"artistIDs"`
	TrackSkipIDs  []*metadata.Gid `json:"trackSkipIDs"`
	ArtistSkipIDs []*metadata.Gid `json:"artistSkipIDs"`
	NumResults    int             `json:"numResults"`
}

type Response struct {
	RecommendedTracks []Track `json:"recommendedTracks"`
	Request           Request `json:"request"`
	Details           Details `json:"details"`
}

// ---

type Track struct {
	Id            string          `json:"id"`
	OriginalId    string          `json:"originalId"`
	Name          string          `json:"name"`
	Artists       []Artist        `json:"artists"`
	Album         Album           `json:"album"`
	Duration      int             `json:"duration"`
	Explicit      bool            `json:"explicit"`
	Popularity    int             `json:"popularity"`
	Score         float64         `json:"score"`
	ContentRating []ContentRating `json:"contentRating"`
}

type Artist struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type Album struct {
	Id            string `json:"id"`
	Name          string `json:"name"`
	LargeImageUrl string `json:"largeImageUrl"`
	ImageUrl      string `json:"imageUrl"`
}

type ContentRating struct {
	Tag     string   `json:"tag"`
	Markets []string `json:"markets"`
}

type Details struct {
	PlaylistLoad float64 `json:"playlistLoad"`
	TotalTime    float64 `json:"totalTime"`
}
