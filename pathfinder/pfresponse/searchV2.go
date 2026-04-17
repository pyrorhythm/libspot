package pfresponse

import (
	pfd "github.com/pyrorhythm/libspot/pathfinder/pfdomain"
)

type SearchV2 struct {
	Albums    *Items[pfd.Album]    `json:"albumsV2,omitempty"`
	Artists   *Items[pfd.Artist]   `json:"artists,omitempty"`
	Episodes  *Items[pfd.Episode]  `json:"episodes,omitempty"`
	Genres    *Items[pfd.Genre]    `json:"genres,omitempty"`
	Playlists *Items[pfd.Playlist] `json:"playlists,omitempty"`
	Podcasts  *Items[pfd.Podcast]  `json:"podcasts,omitempty"`
	Tracks    *Items[pfd.Track]    `json:"tracksV2,omitempty"`
	Users     *Items[pfd.User]     `json:"users,omitempty"`

	TopResults *ItemsV2[pfd.OneofMatched] `json:"topResultsV2,omitempty"`

	Audiobooks *Items[any] `json:"audiobooks,omitempty"`

	ChipOrder pfd.Chips `json:"chipOrder,omitempty"`
	Query     string    `json:"query"`
}

func (s *SearchV2) ToAlbums() *SearchV2Albums {
	return &SearchV2Albums{
		Items: s.Albums,
		Query: s.Query,
	}
}

func (s *SearchV2) ToArtists() *SearchV2Artists {
	return &SearchV2Artists{
		Items: s.Artists,
		Query: s.Query,
	}
}

func (s *SearchV2) ToGenres() *SearchV2Genres {
	return &SearchV2Genres{
		Items: s.Genres,
		Query: s.Query,
	}
}

func (s *SearchV2) ToPlaylists() *SearchV2Playlists {
	return &SearchV2Playlists{
		Items: s.Playlists,
		Query: s.Query,
	}
}

func (s *SearchV2) ToEpisodes() *SearchV2Episodes {
	return &SearchV2Episodes{
		Items: s.Episodes,
		Query: s.Query,
	}
}

func (s *SearchV2) ToPodcasts() *SearchV2Podcasts {
	return &SearchV2Podcasts{
		Items: s.Podcasts,
		Query: s.Query,
	}
}

func (s *SearchV2) ToTracks() *SearchV2Tracks {
	return &SearchV2Tracks{
		Items: s.Tracks,
		Query: s.Query,
	}
}

func (s *SearchV2) ToUsers() *SearchV2Users {
	return &SearchV2Users{
		Items: s.Users,
		Query: s.Query,
	}
}

func (s *SearchV2) ToSuggestions() *SearchV2Suggestions {
	return &SearchV2Suggestions{
		ItemsV2: s.TopResults,
		Query:   s.Query,
	}
}

func (s *SearchV2) ToTop() *SearchV2Top {
	return &SearchV2Top{
		ItemsV2:   s.TopResults,
		Query:     s.Query,
		ChipOrder: s.ChipOrder,
	}
}
