package pathfinder

import (
	"context"
	"fmt"
	"net/http"

	"github.com/cenkalti/backoff/v5"
)

func (p *Pathfinder) makeRequest(
	ctx context.Context,
	req *http.Request,
) (*http.Response, error) {
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/144.0.7559.246 Spotify/1.2.87.414 Safari/537.36")
	req.Header.Set("App-Platform", "OSX_ARM64")
	req.Header.Set("Origin", "https://xpui.app.spotify.com")
	req.Header.Set("Referer", "https://xpui.app.spotify.com")
	req.Header.Set("Content-Type", "application/json;charset=UTF-8")

	return backoff.Retry(
		ctx, func() (*http.Response, error) {
			resp, err := p.client.Do(req.WithContext(ctx))
			if err != nil {
				return nil, backoff.Permanent(err)
			}

			if resp.StatusCode == 401 {
				_ = resp.Body.Close()
				return nil, backoff.Permanent(fmt.Errorf("unauthorized"))
			}

			if resp.StatusCode == 400 {
				_ = resp.Body.Close()
				return nil, backoff.Permanent(fmt.Errorf("bad request"))
			}

			if resp.StatusCode >= 500 {
				_ = resp.Body.Close()
				return nil, backoff.RetryAfter(3)
			}

			return resp, nil
		}, backoff.WithBackOff(backoff.NewExponentialBackOff()))
}
