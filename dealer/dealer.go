package dealer

import (
	"context"
	"errors"
	"fmt"
	"time"

	ws "github.com/coder/websocket"
)

var ErrNotConnected = errors.New("dealer: not connected")

func (d *DealerG2) Start() error {
	if !d.running.CompareAndSwap(false, true) {
		return errors.New("dealer: already started")
	}

	ctx, cancel := context.WithCancel(context.Background())
	d.runningCtx = ctx
	d.runningCancel = cancel

	go d.reconnectLoop(ctx)
	return nil
}

func (d *DealerG2) Stop() error {
	if !d.running.CompareAndSwap(true, false) {
		return nil
	}
	d.runningCancel()
	return nil
}

func (d *DealerG2) Send(msg []byte) error {
	d.mu.RLock()
	ch := d.sendCh
	d.mu.RUnlock()

	if ch == nil {
		return ErrNotConnected
	}

	select {
	case ch <- msg:
		return nil
	default:
		return errors.New("dealer: send buffer full")
	}
}

func (d *DealerG2) reconnectLoop(ctx context.Context) {
	var globalAttempt uint = 0

	for ctx.Err() == nil {
		eps, err := d.fetchEndpoints()
		if err != nil {
			select {
			case <-ctx.Done():
				return
			case <-time.After(time.Duration(d.GlobalRetryDelay) * time.Second):
			}
			continue
		}

		if wsConn, ok := d.tryEndpoints(ctx, eps); ok {
			c := d.newConn(wsConn)
			d.runConn(ctx, c)
			// blocking here
			globalAttempt = 0
			continue
		}

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

func (d *DealerG2) runConn(ctx context.Context, c *conn) {
	d.mu.Lock()
	d.state = ConnStateConnected
	d.sendCh = c.send
	d.mu.Unlock()

	defer func() {
		c.closeWS()
		d.mu.Lock()
		d.sendCh = nil
		d.state = ConnStateIdle
		d.mu.Unlock()
	}()

	c.run(ctx)
}

func (d *DealerG2) newConn(ws *ws.Conn) *conn {
	pingIv := 30 * time.Second
	pingTo := 10 * time.Second
	if d.PingIntervalSec > 0 {
		pingIv = time.Duration(d.PingIntervalSec) * time.Second
	}
	if d.PingTimeout > 0 {
		pingTo = time.Duration(d.PingTimeout) * time.Second
	}

	c := &conn{
		ws:           ws,
		send:         make(chan []byte, 256),
		pingInterval: pingIv,
		pingTimeout:  pingTo,
	}

	c.suicideTimer = time.AfterFunc(pingIv+pingTo, func() { c.Close() })

	return c
}
