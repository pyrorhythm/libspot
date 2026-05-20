package spc

import (
	"context"

	"github.com/cenkalti/backoff/v5"
	"github.com/pkg/errors"
	"github.com/pyrorhythm/fn/bjs"
	"github.com/pyrorhythm/libspot"
	"resty.dev/v3"
)

func makeRequest[to any](
	ctx context.Context,
	rq *resty.Request,
) (*to, error) {
	headers := map[string]string{
		"app-platform": libspot.AppPlatform().String(),
		"accept":       "application/json",
		"user-agent":   "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/144.0.7559.246 Spotify/1.2.87.414 Safari/537.36",
		"origin":       "https://xpui.app.spotify.com",
		"referer":      "https://xpui.app.spotify.com",
		"content-type": "application/json;charset=UTF-8",
	}

	rq = rq.SetHeaders(headers).SetQueryParam("market", "from_token")

	resp, err := backoff.Retry(
		ctx, func() (*resty.Response, error) {
			resp, err := rq.Send()
			if err != nil {
				return nil, backoff.Permanent(err)
			}

			if resp.StatusCode() == 401 {
				return nil, backoff.Permanent(errors.New("unauthorized"))
			}

			if resp.StatusCode() == 400 {
				return nil, backoff.Permanent(errors.New("bad request"))
			}

			if resp.StatusCode() >= 500 {
				return nil, backoff.RetryAfter(3)
			}

			return resp, nil
		}, backoff.WithBackOff(backoff.NewExponentialBackOff()))
	if err != nil {
		return nil, err
	}

	var res *to

	if res, err = bjs.Unmarshal[to](resp.Bytes()); err != nil {
		return nil, errors.Wrapf(err, "failed to unmarshal response")
	}

	return res, nil
}
