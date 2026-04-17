package spc

import (
	"context"
	"encoding/hex"
	"net/http"
	"net/url"

	"github.com/pkg/errors"
	"github.com/pyrorhythm/libspot/spc/metadata"
	"resty.dev/v3"
)

type Typable interface {
	Type() metadata.MdType
}

type Clientable interface {
	Client() *resty.Client
}

func validateGid(id string) error {
	if len(id) != 40 {
		return errors.Errorf("invalid gid %s", id)
	}

	_, err := hex.DecodeString(id)
	if err != nil {
		return errors.Wrapf(err, "invalid gid %s", id)
	}

	return nil
}

func Metadata[T Typable](c *Spclient, ctx context.Context, gid string) (*T, error) {
	var z T

	if err := validateGid(gid); err != nil {
		return nil, err
	}

	u, _ := url.JoinPath(baseUrl, metadata.Path)
	resp, err := makeRequest[T](ctx,
		c.client.R().SetURL(u).
			SetMethod(http.MethodGet).
			SetPathParams(map[string]string{
				"gid":  gid,
				"type": string(z.Type()),
			}),
	)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
