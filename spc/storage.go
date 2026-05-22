package spc

import (
	"context"
	"encoding/hex"
	"fmt"
	"net/http"
	"net/url"

	"github.com/pkg/errors"
	downloadpb "pyrorhythm.dev/libspot/gen/spotify/download"
	metadatapb "pyrorhythm.dev/libspot/gen/spotify/metadata"
)

// ResolveStorageInteractive resolves the CDN storage for an interactive audio
// file in the given format. When prefetch is true it uses the prefetch variant
// of the endpoint, which the client issues to warm the CDN ahead of playback.
func (c *Spclient) ResolveStorageInteractive(
	ctx context.Context,
	fileId []byte,
	format metadatapb.AudioFile_Format,
	prefetch bool,
) (*downloadpb.StorageResolveResponse, error) {
	kind := "interactive"
	if prefetch {
		kind = "interactive_prefetch"
	}

	u, _ := url.JoinPath(baseUrl, fmt.Sprintf(
		"storage-resolve/v2/files/audio/%s/%d/%s",
		kind, format.Number(), hex.EncodeToString(fileId),
	))

	resp := &downloadpb.StorageResolveResponse{}
	if err := makeProtoRequest(ctx,
		c.client.R().SetURL(u).SetMethod(http.MethodGet),
		resp,
	); err != nil {
		return nil, errors.Wrap(err, "failed to resolve interactive storage")
	}

	return resp, nil
}
