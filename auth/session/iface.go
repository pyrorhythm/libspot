package session

import (
	"context"

	"github.com/pyrorhythm/libspot"
)

type Session interface {
	// AuthUrl provides authentication URL to authorize and grant claims.
	AuthUrl(state string) (url, pkce string)
	// AuthCode processes given code with pkce from AuthUrl and authorizes the session.
	AuthCode(ctx context.Context, code, pkce string) error
	// GetOrRefreshToken returns valid (and refreshed if needed) access token.
	GetOrRefreshToken(ctx context.Context) (string, error)
	// GetToken return token without refreshing, even if refresh token is availible.
	GetToken() (string, bool)

	// Load tries to load session from keychain,
	// returning store.Error (store.ErrItemNotFound) if failed
	Load() error
	// Clear clears session from current auth and clears keychain from it (if needed)
	Clear(clearKeychain bool) error

	ClientToken() (string, error)

	Resolver() (libspot.EndpointResolver, error)

	Valid() bool

	RefreshToken() string
	User() string
	DeviceId() string
}
