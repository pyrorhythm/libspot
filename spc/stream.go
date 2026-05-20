package spc

import (
	"context"
	"encoding/hex"
	"net/http"
	"net/url"

	"github.com/cenkalti/backoff/v5"
	"github.com/pkg/errors"
	"github.com/pyrorhythm/libspot"
	downloadpb "github.com/pyrorhythm/libspot/gen/spotify/download"
	metadatapb "github.com/pyrorhythm/libspot/gen/spotify/metadata"
	"google.golang.org/protobuf/proto"
	"resty.dev/v3"
)

// makeProtoRequest sends rq, retrying on transient failures, and unmarshals the
// protobuf response body into msg.
func makeProtoRequest(ctx context.Context, rq *resty.Request, msg proto.Message) error {
	rq = rq.SetHeaders(map[string]string{
		"App-Platform": libspot.AppPlatform().String(),
		"User-Agent":   "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/144.0.7559.246 Spotify/1.2.87.414 Safari/537.36",
		"Accept":       "application/protobuf",
	})

	resp, err := backoff.Retry(ctx, func() (*resty.Response, error) {
		resp, err := rq.Send()
		if err != nil {
			return nil, backoff.Permanent(err)
		}

		switch {
		case resp.StatusCode() == http.StatusUnauthorized:
			return nil, backoff.Permanent(errors.New("unauthorized"))
		case resp.StatusCode() == http.StatusBadRequest:
			return nil, backoff.Permanent(errors.New("bad request"))
		case resp.StatusCode() >= 500:
			return nil, backoff.RetryAfter(3)
		}

		return resp, nil
	}, backoff.WithBackOff(backoff.NewExponentialBackOff()))
	if err != nil {
		return err
	}

	if err = proto.Unmarshal(resp.Bytes(), msg); err != nil {
		return errors.Wrap(err, "failed to unmarshal protobuf response")
	}

	return nil
}

// MetadataTrackProto fetches a track's protobuf metadata. Unlike the JSON
// metadata endpoint, the protobuf variant includes the audio file list
// (file_id + format) needed to resolve and download the stream.
func (c *Spclient) MetadataTrackProto(ctx context.Context, gid string) (*metadatapb.Track, error) {
	if err := validateGid(gid); err != nil {
		return nil, err
	}

	u, _ := url.JoinPath(baseUrl, "metadata/4/track", gid)

	track := &metadatapb.Track{}
	if err := makeProtoRequest(ctx,
		c.client.R().SetURL(u).SetMethod(http.MethodGet),
		track,
	); err != nil {
		return nil, errors.Wrap(err, "failed to fetch track metadata")
	}

	return track, nil
}

// ResolveStorage resolves the CDN URLs for an interactive audio file.
func (c *Spclient) ResolveStorage(ctx context.Context, fileId []byte) ([]string, error) {
	u, _ := url.JoinPath(baseUrl, "storage-resolve/v2/files/audio/interactive/1", hex.EncodeToString(fileId))

	resp := &downloadpb.StorageResolveResponse{}
	if err := makeProtoRequest(ctx,
		c.client.R().SetURL(u).SetMethod(http.MethodGet).SetQueryParam("product", "0"),
		resp,
	); err != nil {
		return nil, errors.Wrap(err, "failed to resolve storage")
	}

	return resp.GetCdnurl(), nil
}
