package resolver

import (
	"context"
	"fmt"
	"net/url"
	"slices"
	"strconv"
	"time"

	"github.com/cenkalti/backoff/v5"
	"github.com/pyrorhythm/fn"
	"github.com/pyrorhythm/fn/bjs"
	"github.com/pyrorhythm/libspot"
	"resty.dev/v3"
)

const resolveEndpoint = "https://apresolve.spotify.com/"

type endpoints struct {
	SpclientEps    []string `json:"spclient"`
	DealerEps      []string `json:"dealer"`
	DealerG2Eps    []string `json:"dealer-g2"`
	AccesspointEps []string `json:"accesspoint"`
}

func (e *endpoints) merge(oth *endpoints) {
	if oth == nil {
		return
	}

	if len(oth.SpclientEps) != 0 {
		e.SpclientEps = oth.SpclientEps
	}

	if len(oth.DealerEps) != 0 {
		e.DealerEps = oth.DealerEps
	}

	if len(oth.DealerG2Eps) != 0 {
		e.DealerG2Eps = oth.DealerG2Eps
	}

	if len(oth.AccesspointEps) != 0 {
		e.AccesspointEps = oth.AccesspointEps
	}
}

func (e endpoints) Spclient() []string {
	return e.SpclientEps
}

func (e endpoints) Dealer() []string {
	return e.DealerEps
}

func (e endpoints) DealerG2() []string {
	return e.DealerG2Eps
}

func (e endpoints) Accesspoint() []string {
	return e.AccesspointEps
}

func tostring(kind libspot.ServiceKind) string {
	return string(kind)
}

type fetcher struct {
	client      *resty.Client
	endpoints   *endpoints
	clientToken string
}

func (t *fetcher) Endpoints() (libspot.Endpoints, bool) {
	return t.endpoints, t.endpoints != nil
}

func (t *fetcher) Fetch(kinds ...libspot.ServiceKind) (libspot.Endpoints, error) {
	v := url.Values{
		"type": fn.Map(kinds, tostring),
	}

	if slices.Contains(kinds, libspot.ServiceKindAccesspoint) {
		v.Set("time", strconv.FormatInt(time.Now().Unix(), 10))
	}

	req := t.client.R().
		SetQueryParamsFromValues(v).
		SetHeaders(map[string]string{
			"client-token": t.clientToken,
			"user-agent":   "Spotify/128600502 (43; 0; 2)",
		})

	_, err := backoff.Retry(context.Background(), func() (any, error) {
		resp, err := req.Get(resolveEndpoint)
		if err != nil {
			// Error is not a bad status but a transport error
			return nil, backoff.Permanent(fmt.Errorf("internal transport failure: %w", err))
		}

		if resp.StatusCode() < 200 || resp.StatusCode() >= 300 {
			return nil, backoff.RetryAfter(3)
		}

		if endp, err := bjs.Unmarshal[endpoints](resp.Bytes()); err != nil {
			return nil, backoff.Permanent(fmt.Errorf("failed to unmarshal endpoints: %w", err))
		} else {
			t.endpoints.merge(endp)
			return nil, nil
		}
	}, backoff.WithBackOff(backoff.NewExponentialBackOff()), backoff.WithMaxTries(5))

	if err != nil {
		return nil, err
	}

	return t.endpoints, nil
}

func New(clientToken string) libspot.EndpointResolver {
	return &fetcher{client: resty.New(), clientToken: clientToken, endpoints: &endpoints{}}
}
