package dealer

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	ws "github.com/coder/websocket"
	"github.com/pyrorhythm/libspot"
)

var ErrEndpointRetriesExceeded = errors.New("dealer: endpoint retries exceeded")

func (d *Dealer) fetchEndpoints() ([]string, error) {
	if eps, ok := d.rslv.Endpoints(); ok {
		if dg2Eps := eps.DealerG2(); len(dg2Eps) > 0 {
			return dg2Eps, nil
		}
	}
	eps, err := d.rslv.Fetch(libspot.ServiceKindDealerG2)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve dealer-g2 endpoints: %w", err)
	}
	dg2Eps := eps.DealerG2()
	if len(dg2Eps) == 0 {
		return nil, fmt.Errorf("no dealer-g2 endpoints resolved")
	}
	return dg2Eps, nil
}

func (d *Dealer) tryEndpoints(ctx context.Context, endpoints []string) (*ws.Conn, bool) {
	for _, ep := range endpoints {
		select {
		case <-ctx.Done():
			return nil, false
		default:
		}
		if wsConn, err := d.tryEndpointWithRetry(ctx, ep); err == nil {
			return wsConn, true
		}
	}
	return nil, false
}

func (d *Dealer) tryEndpointWithRetry(ctx context.Context, endpoint string) (*ws.Conn, error) {
	for attempt := range d.endpoint.Cap {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}

		wsConn, err := d.connectEndpoint(ctx, endpoint)
		if err == nil {
			return wsConn, nil
		}

		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-time.After(d.endpoint.Fn(attempt)):
		}
	}

	return nil, ErrEndpointRetriesExceeded
}

func (d *Dealer) connectEndpoint(ctx context.Context, endpoint string) (*ws.Conn, error) {
	token, err := d.prov.AccessToken(ctx, true)
	if err != nil {
		return nil, fmt.Errorf("token provider: %w", err)
	}

	url := "wss://" + strings.TrimPrefix(endpoint, "https://")

	conn, _, err := ws.Dial(ctx, url, &ws.DialOptions{
		HTTPClient: &http.Client{Timeout: 45 * time.Second},
		HTTPHeader: http.Header{
			"Authorization": {"Bearer " + token},
		},
	})
	if err != nil {
		return nil, fmt.Errorf("dial: %w", err)
	}
	return conn, nil
}
