package auth

import (
	"fmt"

	"github.com/pyrorhythm/libspot"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/spotify"
)

// Spotify OAuth2 scopes.
// See https://developer.spotify.com/documentation/web-api/concepts/scopes.
const (
	ScopeAppRemoteControl          = "app-remote-control"
	ScopePlaylistModify            = "playlist-modify"
	ScopePlaylistModifyPrivate     = "playlist-modify-private"
	ScopePlaylistModifyPublic      = "playlist-modify-public"
	ScopePlaylistRead              = "playlist-read"
	ScopePlaylistReadCollaborative = "playlist-read-collaborative"
	ScopePlaylistReadPrivate       = "playlist-read-private"
	ScopeStreaming                 = "streaming"
	ScopeUGCImageUpload            = "ugc-image-upload"
	ScopeUserFollowModify          = "user-follow-modify"
	ScopeUserFollowRead            = "user-follow-read"
	ScopeUserLibraryModify         = "user-library-modify"
	ScopeUserLibraryRead           = "user-library-read"
	ScopeUserModify                = "user-modify"
	ScopeUserModifyPlaybackState   = "user-modify-playback-state"
	ScopeUserModifyPrivate         = "user-modify-private"
	ScopeUserPersonalized          = "user-personalized"
	ScopeUserReadBirthdate         = "user-read-birthdate"
	ScopeUserReadCurrentlyPlaying  = "user-read-currently-playing"
	ScopeUserReadEmail             = "user-read-email"
	ScopeUserReadPlayHistory       = "user-read-play-history"
	ScopeUserReadPlaybackPosition  = "user-read-playback-position"
	ScopeUserReadPlaybackState     = "user-read-playback-state"
	ScopeUserReadPrivate           = "user-read-private"
	ScopeUserReadRecentlyPlayed    = "user-read-recently-played"
	ScopeUserTopRead               = "user-top-read"
)

func AllScopes() []string {
	return []string{
		ScopeAppRemoteControl,
		ScopePlaylistModify,
		ScopePlaylistModifyPrivate,
		ScopePlaylistModifyPublic,
		ScopePlaylistRead,
		ScopePlaylistReadCollaborative,
		ScopePlaylistReadPrivate,
		ScopeStreaming,
		ScopeUGCImageUpload,
		ScopeUserFollowModify,
		ScopeUserFollowRead,
		ScopeUserLibraryModify,
		ScopeUserLibraryRead,
		ScopeUserModify,
		ScopeUserModifyPlaybackState,
		ScopeUserModifyPrivate,
		ScopeUserPersonalized,
		ScopeUserReadBirthdate,
		ScopeUserReadCurrentlyPlaying,
		ScopeUserReadEmail,
		ScopeUserReadPlayHistory,
		ScopeUserReadPlaybackPosition,
		ScopeUserReadPlaybackState,
		ScopeUserReadPrivate,
		ScopeUserReadRecentlyPlayed,
		ScopeUserTopRead,
	}
}

func NewDefaultOAuthConfig(port int) *oauth2.Config {
	return &oauth2.Config{
		ClientID:    libspot.ClientIdHex,
		RedirectURL: fmt.Sprintf("http://127.0.0.1:%d/login", port),
		Scopes:      AllScopes(),
		Endpoint: oauth2.Endpoint{
			AuthURL:   spotify.Endpoint.AuthURL,
			TokenURL:  spotify.Endpoint.TokenURL,
			AuthStyle: oauth2.AuthStyleInParams,
		},
	}
}
