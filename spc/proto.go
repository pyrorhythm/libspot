package spc

import (
	"context"
	"net/http"

	"github.com/cenkalti/backoff/v5"
	"github.com/pkg/errors"
	"google.golang.org/protobuf/proto"
	"pyrorhythm.dev/libspot"
	"resty.dev/v3"
)

func makeProtoRequest(ctx context.Context, rq *resty.Request, msg proto.Message) error {
	rq = rq.SetContext(ctx).SetHeaders(map[string]string{
		"App-Platform": libspot.AppPlatform().String(),
		"User-Agent":   "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/144.0.7559.246 Spotify/1.2.87.414 Safari/537.36",
		"Accept":       "application/protobuf",
		"Content-Type": "application/protobuf",
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
