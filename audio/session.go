package audio

import (
	"context"
	"fmt"
	"log/slog"
	"math/rand/v2"

	"pyrorhythm.dev/libspot"
	"pyrorhythm.dev/libspot/ap"
	"pyrorhythm.dev/libspot/auth/session"
)

// NewKeyProviderFromSession connects an accesspoint using the session's OAuth
// credentials and returns a KeyProvider ready to request AES audio keys.
func NewKeyProviderFromSession(ctx context.Context, log *slog.Logger, sess session.Session) (*KeyProvider, error) {
	rslv, err := sess.Resolver()
	if err != nil {
		return nil, fmt.Errorf("failed to get resolver from session: %w", err)
	}

	endpoints, err := rslv.Fetch(libspot.ServiceKindAccesspoint)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve accesspoint: %w", err)
	}

	addrs := endpoints.Accesspoint()
	if len(addrs) == 0 {
		return nil, fmt.Errorf("no accesspoint endpoints available")
	}

	addrFn := func(context.Context) string {
		return addrs[rand.IntN(len(addrs))]
	}

	accesspoint := ap.NewAccesspoint(log, addrFn, sess.DeviceId())

	token, err := sess.AccessToken(ctx, true)
	if err != nil {
		return nil, fmt.Errorf("failed to get access token: %w", err)
	}

	if err := accesspoint.ConnectSpotifyToken(ctx, sess.User(), token); err != nil {
		return nil, fmt.Errorf("failed to connect accesspoint: %w", err)
	}

	return NewKeyProvider(log, accesspoint), nil
}
