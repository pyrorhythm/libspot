package pfresponse

import (
	pfd "github.com/pyrorhythm/libspot/pathfinder/pfdomain"
)

type SearchV2 struct {
	Albums     *Items[pfd.Album]          `json:"albumsV2,omitempty"`
	Artists    *Items[pfd.Artist]         `json:"artists,omitempty"`
	Episodes   *Items[pfd.Episode]        `json:"episodes,omitempty"`
	Genres     *Items[pfd.Genre]          `json:"genres,omitempty"`
	Playlists  *Items[pfd.Playlist]       `json:"playlists,omitempty"`
	Podcasts   *Items[pfd.Podcast]        `json:"podcasts,omitempty"`
	Tracks     *Items[pfd.Track]          `json:"tracksV2,omitempty"`
	Users      *Items[pfd.User]           `json:"users,omitempty"`
	TopResults *ItemsV2[pfd.OneofMatched] `json:"topResultsV2,omitempty"`

	Audiobooks *Items[any] `json:"audiobooks,omitempty"`

	ChipOrder pfd.Chips `json:"chipOrder,omitempty"`
	Query     string    `json:"query"`
}
