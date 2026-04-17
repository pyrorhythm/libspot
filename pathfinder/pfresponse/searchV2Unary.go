package pfresponse

import pfd "github.com/pyrorhythm/libspot/pathfinder/pfdomain"

type SearchV2Albums struct {
	*Items[pfd.Album] `json:"albumsV2,omitempty"`

	Query string `json:"query"`
}

type SearchV2Artists struct {
	*Items[pfd.Artist] `json:"artists,omitempty"`

	Query string `json:"query"`
}

type SearchV2Genres struct {
	*Items[pfd.Genre] `json:"genres,omitempty"`

	Query string `json:"query"`
}

type SearchV2Playlists struct {
	*Items[pfd.Playlist] `json:"playlists,omitempty"`

	Query string `json:"query"`
}

type SearchV2Episodes struct {
	*Items[pfd.Episode] `json:"podcasts,omitempty"`

	Query string `json:"query"`
}

type SearchV2Podcasts struct {
	*Items[pfd.Podcast] `json:"podcasts,omitempty"`

	Query string `json:"query"`
}

type SearchV2Tracks struct {
	*Items[pfd.Track] `json:"tracksV2,omitempty"`

	Query string `json:"query"`
}

type SearchV2Users struct {
	*Items[pfd.User] `json:"users,omitempty"`

	Query string `json:"query"`
}

type SearchV2Suggestions struct {
	*ItemsV2[pfd.OneofMatched] `json:"topResultsV2,omitempty"`

	Query string `json:"query"`
}

type SearchV2Top struct {
	*ItemsV2[pfd.OneofMatched] `json:"topResultsV2,omitempty"`

	ChipOrder pfd.Chips `json:"chipOrder,omitempty"`
	Query     string    `json:"query"`
}
