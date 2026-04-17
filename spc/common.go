package spc

import (
	"context"
	"net/http"

	"github.com/goccy/go-json"
	"github.com/pkg/errors"
	"github.com/pyrorhythm/libspot"
	"github.com/pyrorhythm/libspot/auth/client"
	"github.com/pyrorhythm/libspot/spc/extendp"
	"github.com/pyrorhythm/libspot/spc/metadata"
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

func (c *Spclient) MetadataTrack(ctx context.Context, gid string) (*metadata.Track, error) {
	return Metadata[metadata.Track](c, ctx, gid)
}

func (c *Spclient) MetadataAlbum(ctx context.Context, gid string) (*metadata.Album, error) {
	return Metadata[metadata.Album](c, ctx, gid)
}

func (c *Spclient) MetadataArtist(ctx context.Context, gid string) (*metadata.Artist, error) {
	return Metadata[metadata.Artist](c, ctx, gid)
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
