package session

import "context"

type Session interface {
	// AuthUrl provides authentication URL to authorize and grant claims.
	AuthUrl(state string) (url, pkce string)
	// AuthCode processes given code with pkce from AuthUrl and authorizes the session.
	AuthCode(ctx context.Context, code, pkce string) error
	// AccessToken returns valid (and refreshed if needed) access token.
	AccessToken(ctx context.Context) (string, error)

	// Load tries to load session from keychain,
	// returning store.Error (store.ErrItemNotFound) if failed
	Load() error
	// Clear clears session from current auth and clears keychain from it (if needed)
	Clear(clearKeychain bool) error

	Valid() bool

	RefreshToken() string
	User() string
	DeviceId() string
}
