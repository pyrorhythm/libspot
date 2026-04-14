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
				return nil, backoff.Permanent(fmt.Errorf("bad request"))
			}

			if resp.StatusCode >= 500 {
				return nil, backoff.RetryAfter(3)
			}

			return resp, nil
		}, backoff.WithBackOff(backoff.NewExponentialBackOff()))
}
