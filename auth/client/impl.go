package client

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/pyrorhythm/libspot"
	"github.com/pyrorhythm/zlog"
)

type transport struct {
	transportSettings

	prov libspot.TokenProvider
}

func (t *transport) RoundTrip(req *http.Request) (*http.Response, error) {
	ctx := req.Context()
	ctx = zlog.AddToContext(
		ctx,
		slog.String("url", req.URL.String()),
		slog.String("method", req.Method),
	)

	req = req.Clone(ctx)
	if t.injectClientToken {
		slog.Log(ctx, zlog.LevelDebug, "injecting client token")

		ctok, err := t.prov.ClientToken()
		if err != nil {
			slog.Log(ctx, zlog.LevelError, "failed to inject client token", "error", err)
			return nil, fmt.Errorf("client: failed to get client token: %w", err)
		}
		ctx = zlog.AddToContext(ctx, slog.String("clientToken", ctok))

		slog.Log(ctx, zlog.LevelDebug, "got client token")
		req.Header.Set("Client-Token", ctok)
	}
	if t.injectAccessToken {
		slog.Log(ctx, zlog.LevelDebug, "injecting access token")

		act, err := t.prov.AccessToken(req.Context(), t.canRefreshAccessToken)
		if err != nil {
			slog.Log(ctx, zlog.LevelError, "failed to inject access token", "error", err)

			return nil, fmt.Errorf(
				"client: failed to get access token (triedToRefresh=%v): %w",
				t.canRefreshAccessToken, err,
			)
		}
		ctx = zlog.AddToContext(ctx, slog.String("accessToken", act))
		slog.Log(ctx, zlog.LevelDebug, "got access token")
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", act))
	}

	slog.Log(ctx, zlog.LevelDebug, "sending request", "req", req)
	return t.baseTransport.RoundTrip(req)
}
