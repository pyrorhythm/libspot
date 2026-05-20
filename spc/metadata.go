package spc

import (
	"context"
	"encoding/hex"
	"net/http"

	"github.com/pkg/errors"
	"github.com/pyrorhythm/libspot/spc/metadata"
	"resty.dev/v3"
)

const origin = "https://xpui.app.spotify.com"

// preflight mirrors the CORS preflight the Spotify web client (xpui) issues
// before the metadata GET. The server does not strictly require it for a
// non-browser client, but we replicate the real request flow.
func preflight(ctx context.Context, rq *resty.Request, method, requestHeaders string) error {
	resp, err := rq.
		SetContext(ctx).
		SetMethod(http.MethodOptions).
		SetHeaders(map[string]string{
			"Accept":                         "*/*",
			"Access-Control-Request-Method":  method,
			"Access-Control-Request-Headers": requestHeaders,
			"Origin":                         origin,
			"Referer":                        origin + "/",
			"Sec-Fetch-Mode":                 "cors",
			"Sec-Fetch-Site":                 "same-site",
			"Sec-Fetch-Dest":                 "empty",
		}).
		Send()
	if err != nil {
		return errors.Wrap(err, "preflight request failed")
	}
	if resp.StatusCode() >= 400 {
		return errors.Errorf("preflight rejected with status %d", resp.StatusCode())
	}
	return nil
}

type Typable interface {
	Type() metadata.MdType
}

func validateGid(id string) (err error) {
	if _, err = hex.DecodeString(id); err != nil {
		err = errors.Wrapf(err, "invalid gid %s", id)
	}
	return err
}

func Metadata[T Typable](c *Spclient, ctx context.Context, gid string) (*T, error) {
	var z T

	if err := validateGid(gid); err != nil {
		return nil, err
	}

	u := baseUrl + metadata.Path
	pathParams := map[string]string{
		"gid":  gid,
		"type": string(z.Type()),
	}

	if err := preflight(ctx,
		c.client.R().SetURL(u).SetPathParams(pathParams),
		http.MethodGet,
		"app-platform,authorization,client-token,spotify-app-version",
	); err != nil {
		return nil, err
	}

	resp, err := makeRequest[T](ctx,
		c.client.R().SetURL(u).
			SetMethod(http.MethodGet).
			SetPathParams(pathParams),
	)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
