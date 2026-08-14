package pfrequest

import (
	pfd "pyrorhythm.dev/libspot/pathfinder/pfdomain"
)

type BadgeOperation Operation

const (
	OpSearchTracks      BadgeOperation = "searchTracks"
	OpSearchAlbums      BadgeOperation = "searchAlbums"
	OpSearchArtists     BadgeOperation = "searchArtists"
	OpSearchPlaylists   BadgeOperation = "searchPlaylists"
	OpSearchPodcasts    BadgeOperation = "searchPodcasts"
	OpSearchEpisodes    BadgeOperation = "searchFullEpisodes"
	OpSearchUsers       BadgeOperation = "searchUsers"
	OpSearchGenres      BadgeOperation = "searchGenres"
	OpSearchDesktop     Operation      = "searchDesktop"
	OpSearchTop         Operation      = "searchTopResultsList"
	OpSearchSuggestions Operation      = "searchSuggestions"

	OpHome                   Operation = "home"
	OpGetAlbum               Operation = "getAlbum"
	OpWhatsNewFeedNewItems   Operation = "whatsNewFeedNewItems"
	OpRecentSearches         Operation = "recentSearches"
	OpSaveRecentSearches     Operation = "saveRecentSearches"
	OpIsFollowingUsers       Operation = "isFollowingUsers"
	OpFeedBaselineLookup     Operation = "feedBaselineLookup"
	OpFetchExtractedColors   Operation = "fetchExtractedColors"
	OpGetDynamicColorsByUris Operation = "getDynamicColorsByUris"
	OpQueryNpvArtist         Operation = "queryNpvArtist"
)

const (
	opSearchTracks    = Operation(OpSearchTracks)
	opSearchAlbums    = Operation(OpSearchAlbums)
	opSearchArtists   = Operation(OpSearchArtists)
	opSearchPlaylists = Operation(OpSearchPlaylists)
	opSearchPodcasts  = Operation(OpSearchPodcasts)
	opSearchEpisodes  = Operation(OpSearchEpisodes)
	opSearchUsers     = Operation(OpSearchUsers)
	opSearchGenres    = Operation(OpSearchGenres)
)

func (b BadgeOperation) String() string {
	return string(b)
}

func (b BadgeOperation) Valid() bool {
	switch b {
	case OpSearchTracks, OpSearchAlbums, OpSearchArtists, OpSearchPlaylists, OpSearchUsers:
		return true
	}

	return false
}

type Operation string

var generatedOperationHashes map[Operation]string

func (o Operation) graphQLHash() string {
	if hash, ok := generatedOperationHashes[o]; ok {
		return hash
	}

	panic("not implemented")
}

func (o Operation) Extension() *pfd.PersistedQuery {
	return &pfd.PersistedQuery{
		Version:    1,
		Sha256Hash: o.graphQLHash(),
	}
}
