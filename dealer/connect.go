package dealer

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	ws "github.com/coder/websocket"
	"github.com/pyrorhythm/libspot"
)

func (d *DealerG2) fetchEndpoints() ([]string, error) {
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

func (d *DealerG2) tryEndpoints(ctx context.Context, endpoints []string) (*ws.Conn, bool) {
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

func (d *DealerG2) tryEndpointWithRetry(ctx context.Context, endpoint string) (*ws.Conn, error) {
	for attempt := uint(0); ; attempt++ {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}

		wsConn, err := d.connectEndpoint(ctx, endpoint)
		if err == nil {
			return wsConn, nil
		}

		delay := d.delay(time.Second, attempt)
		t := time.NewTimer(delay)
		select {
		case <-ctx.Done():
			t.Stop()
			return nil, ctx.Err()
		case <-t.C:
		}
	}
}

func (d *DealerG2) connectEndpoint(ctx context.Context, endpoint string) (*ws.Conn, error) {
	token, err := d.prov.AccessToken(ctx)
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
