package pfrequest

import (
	pfd "github.com/pyrorhythm/libspot/pathfinder/pfdomain"
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

func (o Operation) graphQLHash() string {
	switch o {
	case OpHome:
		return "23e37f2e58d82d567f27080101d36609009d8c3676457b1086cb0acc55b72a5d"
	case OpFeedBaselineLookup:
		return "a950fb7c4ecdcaf2aad2f3ca9ee9c3aa4b9c43c97e1d07d05148c4d355bea7fc"
	case OpFetchExtractedColors:
		return "36e90fcaea00d47c695fce31874efeb2519b97d4cd0ee1abfb4f8dc9348596ea"
	case OpGetDynamicColorsByUris:
		return "f0f112945d6d745bd8ff790317bbf8d310036da75df33130490e9d6dc96c59d9"
	case OpSaveRecentSearches, OpRecentSearches:
		return "873c63d8337500d59512c3dc80f47d1b1ab4f12c94fae8136c6a3103e1c976a5"
	case OpSearchSuggestions:
		return "1b44e7bced744d15c47e6c4c11952541693324020c528dc97d19c4a38cfb754e"
	case OpWhatsNewFeedNewItems:
		return "d889c8c936ab192af8ced595427f5ba2acdf63478fdc0a181c8d477f8322630e"
	case OpIsFollowingUsers:
		return "c00e0cb6c7766e7230fc256cf4fe07aec63b53d1160a323940fce7b664e95596"
	case OpSearchTop:
		return "43314b043ad59fe5d06d6a812be70ab3c062ee72633ad9ee460f0e74a86ef7c5"
	case opSearchUsers:
		return "d3f7547835dc86a4fdf3997e0f79314e7580eaf4aaf2f4cb1e71e189c5dfcb1f"
	case opSearchEpisodes:
		return "02f66233401a3cd965c038a2c2fede2911dd93517699b868cb3239e7e17a64f5"
	case opSearchPodcasts:
		return "0195d9f61b43606d490bca64c3456e3593528cea6cc05c7e822c7c42beed0f4e"
	case opSearchPlaylists:
		return "af1730623dc1248b75a61a18bad1f47f1fc7eff802fb0676683de88815c958d8"
	case opSearchAlbums:
		return "5e7d2724fbef31a25f714844bf1313ffc748ebd4bd199eaad50628a4f246a7ab"
	case opSearchArtists:
		return "72c8c7c1e789a9f11e261c4f9ae35a9465bbb90137c584428989573617b6c08d"
	case opSearchTracks:
		return "59ee4a659c32e9ad894a71308207594a65ba67bb6b632b183abe97303a51fa55"
	case opSearchGenres:
		return "9e1c0e056c46239dd1956ea915b988913c87c04ce3dadccdb537774490266f46"
	}

	panic("not implemented")
}

func (o Operation) Extension() *pfd.PersistedQuery {
	return &pfd.PersistedQuery{
		Version:    1,
		Sha256Hash: o.graphQLHash(),
	}
}
