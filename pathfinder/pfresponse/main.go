package pfresponse

import pfd "github.com/pyrorhythm/libspot/pathfinder/pfdomain"

type Response struct {
	SearchV2       *SearchV2       `json:"searchV2"`
	RecentSearches *RecentSearches `json:"recentSearches"`
	AlbumUnion     *pfd.AlbumFull  `json:"albumUnion"`

	DynamicColors   []*pfd.DynamicColorSet        `json:"dynamicColors"`
	Lookup          []*pfd.PlaylistPreviewItems   `json:"lookup"`
	IsFollowing     []*pfd.IsFollowingUser        `json:"users"`
	ExtractedColors []*pfd.ExtractedColorExtended `json:"extractedColors"`
}

type Payload[T any] struct {
	Data       *T              `json:"data"`
	Extensions *pfd.Extensions `json:"extensions,omitempty"`
}

func (p *Payload[T]) Get() *T {
	return p.Data
}

type Items[T any] struct {
	pfd.ItemList[T]
	pfd.TotalCount
	PagingInfo `json:"pagingInfo"`
}

type ItemsV2[T any] struct {
	pfd.ItemV2List[T]
	pfd.TotalCount
	PagingInfo `json:"pagingInfo"`
}

type PagingInfo struct {
	Limit      int `json:"limit"`
	NextOffset int `json:"nextOffset"`
}

func (p PagingInfo) GetLimit() int {
	return p.Limit
}

func (p PagingInfo) GetNextOffset() int {
	return p.NextOffset
}

