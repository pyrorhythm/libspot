package spc

import (
	"context"
	"net/http"

	"github.com/goccy/go-json"
	"github.com/pkg/errors"
	"github.com/pyrorhythm/libspot"
	"github.com/pyrorhythm/libspot/auth/client"
	"github.com/pyrorhythm/libspot/spc/extendp"
	"resty.dev/v3"
)

const baseUrl = "https://spclient.wg.spotify.com/"

type Spclient struct {
	client *resty.Client
	prov   libspot.TokenProvider
	endp   libspot.EndpointResolver
}

func New(prov libspot.TokenProvider, endp libspot.EndpointResolver) *Spclient {
	return &Spclient{
		prov: prov,
		endp: endp,
		client: client.
			NewAuthorizedClient(prov, client.CanRefreshAccessToken(true)).
			Client(),
	}
}

func (c *Spclient) ExtendPlaylist(
	ctx context.Context,
	req *extendp.Request,
) (*extendp.Response, error) {
	bs, err := json.Marshal(req)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to marshal extend playlist request")
	}

	resp, err := makeRequest[extendp.Response](
		ctx, c.client.R().SetMethod(http.MethodPost).SetBody(bs),
	)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to send extend playlist request")
	}

	return resp, nil
}
