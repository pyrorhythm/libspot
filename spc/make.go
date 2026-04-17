package spc

import (
	"context"
	"fmt"

	"github.com/cenkalti/backoff/v5"
	"github.com/goccy/go-json"
	"github.com/pkg/errors"
	"github.com/pyrorhythm/libspot"
	"resty.dev/v3"
)

func makeRequest[to any](
	ctx context.Context,
	rq *resty.Request,
) (*to, error) {
	headers := map[string]string{
		"App-Platform": libspot.AppPlatform().String(),
		"User-Agent":   "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/144.0.7559.246 Spotify/1.2.87.414 Safari/537.36",
		"Origin":       "https://xpui.app.spotify.com",
		"Referer":      "https://xpui.app.spotify.com",
		"Content-Type": "application/json;charset=UTF-8",
	}

	rq = rq.SetHeaders(headers)

	resp, err := backoff.Retry(
		ctx, func() (*resty.Response, error) {
			resp, err := rq.Send()
			if err != nil {
				return nil, backoff.Permanent(err)
			}

			if resp.StatusCode() == 401 {
				return nil, backoff.Permanent(fmt.Errorf("unauthorized"))
			}

			if resp.StatusCode() == 400 {
				return nil, backoff.Permanent(fmt.Errorf("bad request"))
			}

			if resp.StatusCode() >= 500 {
				return nil, backoff.RetryAfter(3)
			}

			return resp, nil
		}, backoff.WithBackOff(backoff.NewExponentialBackOff()))
	if err != nil {
		return nil, err
	}

	var res to

	if err = json.Unmarshal(resp.Bytes(), &res); err != nil {
		return nil, errors.Wrapf(err, "failed to unmarshal response")
	}

	return &res, nil
}
