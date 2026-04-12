package dealer

import (
	"context"
	"fmt"
	"time"

	ws "github.com/coder/websocket"
	"github.com/pyrorhythm/libspot/dealer/types"
)

func (d *Dealer) loop(ctx context.Context) {
	var globalAttempt int64

	for {
		select {
		case <-ctx.Done():
			return
		default:
		}

		eps, err := d.fetchEndpoints()
		if err != nil {
			select {
			case <-ctx.Done():
				return
			case <-time.After(time.Duration(d.GlobalRetryDelay) * time.Second):
			}
			continue
		}

		d.connMu.Lock()
		if wsConn, ok := d.tryEndpoints(ctx, eps); ok {
			d.newConn(wsConn)
			d.connMu.Unlock()
			// /
			d.runConn(ctx)
			// blocking here
			globalAttempt = 0
			continue
		}
		d.connMu.Unlock()

		globalAttempt++
		if d.RetryCap > 0 && globalAttempt > d.RetryCap {
			panic(
				fmt.Sprintf("dealer: all endpoints exhausted after %d attempts", globalAttempt),
			)
		}

		select {
		case <-ctx.Done():
			return
		case <-time.After(time.Duration(d.GlobalRetryDelay) * time.Second):
		}
	}
}

func (d *Dealer) runConn(ctx context.Context) {
	defer func() {
		d.connMu.Lock()
		defer d.connMu.Unlock()

		d.conn.closeWS()
	}()

	d.conn.run(ctx)
}

func (d *Dealer) newConn(ws *ws.Conn) {
	pingIv := 30 * time.Second
	pingTo := 10 * time.Second
	if d.PingIntervalSec > 0 {
		pingIv = time.Duration(d.PingIntervalSec) * time.Second
	}
	if d.PingTimeout > 0 {
		pingTo = time.Duration(d.PingTimeout) * time.Second
	}

	d.conn = &conn{
		dealer:       d,
		ws:           ws,
		send:         make(chan []byte, 256),
		reqCh:        make(chan *types.Request, 64),
		pingInterval: pingIv,
		pingTimeout:  pingTo,
	}

	d.conn.suicideTimer = time.AfterFunc(pingIv+pingTo, func() { _ = d.conn.Close() })
}
