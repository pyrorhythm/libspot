package pfresponse

import (
	"github.com/pyrorhythm/libspot/pathfinder/pfdomain"
)

type PagingInfo struct {
	Limit      int `json:"limit"`
	NextOffset int `json:"nextOffset"`
}

type TopResultsItems struct {
	Items []TopResultHit `json:"itemsV2"`
	// Those are missing on OpTopSearch
	PagingInfo *PagingInfo `json:"pagingInfo,omitempty"`
	TotalCount *int        `json:"totalCount,omitempty"`
}

type Counts struct {
	AlbumsCount     *pfdomain.TotalCount `json:"albumsV2,omitempty"`
	ArtistsCount    *pfdomain.TotalCount `json:"artists,omitempty"`
	AudiobooksCount *pfdomain.TotalCount `json:"audiobooks,omitempty"`
	EpisodesCount   *pfdomain.TotalCount `json:"episodes,omitempty"`
	GenresCount     *pfdomain.TotalCount `json:"genres,omitempty"`
	PlaylistsCount  *pfdomain.TotalCount `json:"playlists,omitempty"`
	PodcastsCount   *pfdomain.TotalCount `json:"podcasts,omitempty"`
	TracksCount     *pfdomain.TotalCount `json:"tracksV2,omitempty"`
	UsersCount      *pfdomain.TotalCount `json:"users,omitempty"`
}

type SearchResultV2 struct {
	SearchResultV2Payload `json:"searchV2"`
}

type SearchResultV2Payload struct {
	// [Counts] is availible only for OpTopSearch
	Counts

	ChipOrder pfdomain.Chips  `json:"chipOrder,omitempty"`
	Query     string          `json:"query"`
	Results   TopResultsItems `json:"topResultsV2"`
}
