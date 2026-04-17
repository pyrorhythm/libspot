package pfresponse

import pfd "github.com/pyrorhythm/libspot/pathfinder/pfdomain"

type RecentSearches struct {
	Items pfd.ItemList[pfd.Oneof] `json:"recentSearchesItems"`
}
