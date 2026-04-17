package pfresponse

import pfd "github.com/pyrorhythm/libspot/pathfinder/pfdomain"

type (
	DynamicColors   []*pfd.DynamicColorSet
	ExtractedColors []*pfd.ExtractedColorExtended
	Lookup          []*pfd.PlaylistPreviewItems
	IsFollowing     []*pfd.IsFollowingUser
	RecentSearches  struct {
		Items pfd.ItemList[pfd.Oneof] `json:"recentSearchesItems"`
	}
)
