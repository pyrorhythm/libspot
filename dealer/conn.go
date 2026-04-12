package dealer

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	ws "github.com/coder/websocket"
	"golang.org/x/sync/errgroup"
)

const maxReqHandlers = 16

type conn struct {
	dealer *DealerG2 // backref to access router

	ws           *ws.Conn
	send         chan []byte
	reqCh        chan *Msg // buffered channel for incoming requests
	pingInterval time.Duration
	pingTimeout  time.Duration

	suicideTimer *time.Timer

	// coordination
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

		var m Msg
		if err := json.Unmarshal(msg, &m); err != nil {
			continue
		}

		if m.Type == "request" {
			select {
			case c.reqCh <- &m:
			default:
			}
			continue
		}

		switch m.Type {
		case "message":
			c.dealer.router.handleMsg(&m)
		default:
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

func (c *conn) handleReq(msg *Msg) {
	replyCh, found := c.dealer.router.handleReq(c.ctx, msg)
	if !found {
		return
	}

	var ok bool
	select {
	case <-c.ctx.Done():
		return
	case ok = <-replyCh:
	}

	reply := Reply{Type: "response", Key: msg.Key}
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
