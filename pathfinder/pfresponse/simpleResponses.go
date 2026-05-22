package pfresponse

import pfd "pyrorhythm.dev/libspot/pathfinder/pfdomain"

type RecentSearches struct {
	Items pfd.ItemList[pfd.Oneof] `json:"recentSearchesItems"`
}
