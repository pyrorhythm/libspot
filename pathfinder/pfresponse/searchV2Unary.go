package pfresponse

import pfd "github.com/pyrorhythm/libspot/pathfinder/pfdomain"

type SearchV2Albums struct {
	Albums *Items[pfd.Album] `json:"albumsV2,omitempty"`
	Query  string            `json:"query"`
}

type SearchV2Artists struct {
	Artists *Items[pfd.Artist] `json:"artists,omitempty"`
	Query   string             `json:"query"`
}

type SearchV2Genres struct {
	Genres *Items[pfd.Genre] `json:"genres,omitempty"`
	Query  string            `json:"query"`
}

type SearchV2Playlists struct {
	Playlists *Items[pfd.Playlist] `json:"playlists,omitempty"`
	Query     string               `json:"query"`
}

type SearchV2Podcasts struct {
	Podcasts *Items[pfd.Podcast] `json:"podcasts,omitempty"`
	Query    string              `json:"query"`
}

type SearchV2Tracks struct {
	Tracks *Items[pfd.Track] `json:"tracksV2,omitempty"`
	Query  string            `json:"query"`
}

type SearchV2Users struct {
	Users *Items[pfd.User] `json:"users,omitempty"`
	Query string           `json:"query"`
}

type SearchV2Suggestions struct {
	TopResults *ItemsV2[pfd.OneofMatched] `json:"topResultsV2,omitempty"`
	Query      string                     `json:"query"`
}

type SearchV2Top struct {
	TopResults *ItemsV2[pfd.OneofMatched] `json:"topResultsV2,omitempty"`
	ChipOrder  pfd.Chips                  `json:"chipOrder,omitempty"`
	Query      string                     `json:"query"`
}
