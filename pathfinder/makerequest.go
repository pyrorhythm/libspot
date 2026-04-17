package pathfinder

import (
	"context"
	"fmt"

	"github.com/cenkalti/backoff/v5"
	"github.com/pyrorhythm/libspot"
	pfq "github.com/pyrorhythm/libspot/pathfinder/pfrequest"
	pfs "github.com/pyrorhythm/libspot/pathfinder/pfresponse"
	"resty.dev/v3"
)

func (p *Pathfinder) makeRequest(
	ctx context.Context,
	rq pfq.Request,
) (*pfs.Response, error) {
	body, err := marshal(rq)

	if err != nil {
		return nil, err
	}

	headers := map[string]string{
		"App-Platform": libspot.AppPlatform().String(),
		"User-Agent":   "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/144.0.7559.246 Spotify/1.2.87.414 Safari/537.36",
		"Origin":       "https://xpui.app.spotify.com",
		"Referer":      "https://xpui.app.spotify.com",
		"Content-Type": "application/json;charset=UTF-8",
	}

	resp, err := backoff.Retry(
		ctx, func() (*resty.Response, error) {
			resp, err := p.client.R().
				SetContext(ctx).
				SetHeaders(headers).
				SetBody(body).
				Post(reqUrl)
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

	pl, err := unmarshal(resp.Bytes())

	if err != nil {
		return nil, err
	}

	return pl.Get(), nil
}
