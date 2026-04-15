package responsetypes

import (
	"github.com/pyrorhythm/libspot/pathfinder/domaintypes"
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
	AlbumsCount     *domaintypes.TotalCount `json:"albumsV2,omitempty"`
	ArtistsCount    *domaintypes.TotalCount `json:"artists,omitempty"`
	AudiobooksCount *domaintypes.TotalCount `json:"audiobooks,omitempty"`
	EpisodesCount   *domaintypes.TotalCount `json:"episodes,omitempty"`
	GenresCount     *domaintypes.TotalCount `json:"genres,omitempty"`
	PlaylistsCount  *domaintypes.TotalCount `json:"playlists,omitempty"`
	PodcastsCount   *domaintypes.TotalCount `json:"podcasts,omitempty"`
	TracksCount     *domaintypes.TotalCount `json:"tracksV2,omitempty"`
	UsersCount      *domaintypes.TotalCount `json:"users,omitempty"`
}

type SearchResultV2 struct {
	SearchResultV2Payload `json:"searchV2"`
}

type SearchResultV2Payload struct {
	// [Counts] is availible only for OpTopSearch
	Counts

	ChipOrder domaintypes.Chips `json:"chipOrder,omitempty"`
	Query     string            `json:"query"`
	Results   TopResultsItems   `json:"topResultsV2"`
}
