package dealer

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"time"

	ws "github.com/coder/websocket"
	"github.com/pyrorhythm/libspot/dealer/types"
	"golang.org/x/sync/errgroup"
)

const maxReqHandlers = 16

type conn struct {
	dealer *Dealer // backref to access router

	ws           *ws.Conn
	send         chan []byte
	reqCh        chan *types.Request
	pingInterval time.Duration
	pingTimeout  time.Duration

	suicideTimer *time.Timer

	// lifespan
	ctx    context.Context
	cancel context.CancelFunc
	wg     *errgroup.Group
}

func (c *conn) run(ctx context.Context) {
	c.ctx, c.cancel = context.WithCancel(ctx)
	defer c.cancel()

	c.wg, c.ctx = errgroup.WithContext(c.ctx)

	c.wg.Go(c.sendLoop)
	c.wg.Go(c.pingLoop)
	c.wg.Go(c.recvLoop)
	c.wg.Go(c.reqLoop)

	err := c.wg.Wait()
	c.closeWS()

	if c.suicideTimer != nil {
		c.suicideTimer.Stop()
	}

	if c.dealer.OnClose != nil {
		c.dealer.OnClose(err)
	}
}

func (c *conn) recvLoop() error {
	for {
		select {
		case <-c.ctx.Done():
			return c.ctx.Err()
		default:
		}

		_, msg, err := c.ws.Read(c.ctx)
		if err != nil {
			wsErr, ok := errors.AsType[ws.CloseError](err)

			if ok && wsErr.Code == ws.StatusNormalClosure {
				return nil
			}
			return err
		}

		c.resetSuicideTimer()

		var env types.Envelope
		if err := json.Unmarshal(msg, &env); err != nil {
			continue
		}

		slog.Debug("[dealer] recv envelope", "typ", env.Type, "uri", env.Uri)

		switch {
		case env.IsMessage():
			c.dealer.router.handleMsg(env.ToMessage())
		case env.IsRequest():
			req, err := env.ToRequest()
			if err != nil {
				continue
			}
			select {
			case c.reqCh <- req:
			default:
			}
		}
	}
}

func (c *conn) reqLoop() error {
	sem := make(chan struct{}, maxReqHandlers)

	for {
		select {
		case <-c.ctx.Done():
			return nil
		case msg := <-c.reqCh:
			sem <- struct{}{}
			c.wg.Go(func() error {
				defer func() { <-sem }()
				c.handleReq(msg)
				return nil
			})
		}
	}
}

func (c *conn) handleReq(req *types.Request) {
	replyCh, found := c.dealer.router.handleReq(c.ctx, req)
	if !found {
		return
	}

	var ok bool
	select {
	case <-c.ctx.Done():
		return
	case ok = <-replyCh:
	}

	reply := types.Reply{Type: "response", Key: req.Key}
	reply.Payload.Success = ok

	payload, err := json.Marshal(reply)
	if err != nil {
		return
	}

	select {
	case c.send <- payload:
	default: // send full, reply dropped
	}
}

func (c *conn) sendLoop() error {
	for {
		select {
		case <-c.ctx.Done():
			return nil
		case msg, ok := <-c.send:
			if !ok {
				return errors.New("channel closed")
			}
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			if err := c.ws.Write(ctx, ws.MessageText, msg); err != nil {
				cancel()
				return err
			}
			cancel()
		}
	}
}

func (c *conn) pingLoop() error {
	ticker := time.NewTicker(c.pingInterval)
	defer ticker.Stop()

	for {
		select {
		case <-c.ctx.Done():
			return nil
		case <-ticker.C:
			ctx, cancel := context.WithTimeout(context.Background(), c.pingTimeout)
			if err := c.ws.Ping(ctx); err != nil {
				cancel()
				return err
			}
			cancel()
		}
	}
}

func (c *conn) resetSuicideTimer() {
	if c.suicideTimer == nil {
		return
	}
	if !c.suicideTimer.Stop() { // already fired, ws conn is dead
		return
	}
	c.suicideTimer.Reset(c.pingInterval + c.pingTimeout)
}

func (c *conn) Close() error {
	if c.cancel != nil {
		c.cancel()
	}
	return nil
}

func (c *conn) closeWS() {
	c.ws.CloseNow()
}
