package dealer

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"pyrorhythm.dev/fn"
	"pyrorhythm.dev/libspot"
	"pyrorhythm.dev/libspot/auth/session"
	"pyrorhythm.dev/libspot/dealer/types"
	"pyrorhythm.dev/libspot/pkg/delay"
)

var (
	ErrNotConnected = errors.New("dealer: not connected")
	ErrSendOverflow = errors.New("dealer: send buffer full")
)

type DelayConfig struct {
	Fn  delay.Delay
	Cap int64
}

type Dealer struct {
	prov  libspot.TokenProvider
	rslv  libspot.EndpointResolver
	delay func(attempt int64) time.Duration

	endpoint *DelayConfig
	global   *DelayConfig

	interval time.Duration
	timeout  time.Duration

	onConnectionShutdown func(error)

	router Router

	conn   *conn
	connMu sync.RWMutex

	cidMu        sync.RWMutex
	connectionId string

	running    atomic.Bool
	cancelLoop context.CancelFunc
	done       chan struct{}
}

var commonDelayCfg = &DelayConfig{
	Fn:  delay.ExponentialJitter2Delay(2 * time.Second),
	Cap: 5,
}

func applyDefaults(d *Dealer) *Dealer {
	d.endpoint = commonDelayCfg
	d.global = commonDelayCfg

	d.interval = 10 * time.Second
	d.timeout = 30 * time.Second

	d.onConnectionShutdown = func(error) {}

	return d
}

func coverNils(dealer *Dealer) {
	if dealer.endpoint == nil && dealer.global == nil {
		dealer.endpoint = commonDelayCfg
		dealer.global = commonDelayCfg
	} else if dealer.endpoint == nil {
		dealer.endpoint = dealer.global
	} else if dealer.global == nil {
		dealer.global = commonDelayCfg
	}

	if dealer.interval <= 0 {
		dealer.interval = 10 * time.Second
	}

	if dealer.timeout <= 0 {
		dealer.timeout = 30 * time.Second
	}

	if dealer.onConnectionShutdown == nil {
		dealer.onConnectionShutdown = func(error) {}
	}
}

func New(
	prov libspot.TokenProvider,
	rslv libspot.EndpointResolver,
	opts ...Option,
) *Dealer {
	d := &Dealer{
		prov:   prov,
		rslv:   rslv,
		router: newRouter(),
	}

	d = fn.Apply(applyDefaults(d), opts...)

	if d.endpoint == nil && d.global == nil {
		d.endpoint = commonDelayCfg
		d.global = commonDelayCfg
	} else if d.endpoint == nil {
		d.endpoint = d.global
	} else if d.global == nil {
		d.global = commonDelayCfg
	}
	if d.interval <= 0 {
		d.interval = 10 * time.Second
	}
	if d.timeout <= 0 {
		d.timeout = 30 * time.Second
	}
	if d.onConnectionShutdown == nil {
		d.onConnectionShutdown = func(error) {}
	}

	return d
}

func NewFromSession(
	sess session.Session,
	opts ...Option,
) (*Dealer, error) {
	rslv, err := sess.Resolver()
	if err != nil {
		return nil, fmt.Errorf("failed to get resolver from session: %w", err)
	}
	return New(sess, rslv, opts...), nil
}

// ConnectionID returns the pusher connection id assigned to the live dealer
// socket, or "" before the first connection frame arrives. Connect-state
// controllers register their hidden device against this id so Spotify pushes
// player state back over this same socket.
func (d *Dealer) ConnectionID() string {
	d.cidMu.RLock()
	defer d.cidMu.RUnlock()
	return d.connectionId
}

func (d *Dealer) setConnectionID(id string) {
	d.cidMu.Lock()
	defer d.cidMu.Unlock()
	d.connectionId = id
}

func (d *Dealer) OnMsg(uri string, cb func(*types.Message)) (unsubscribe func()) {
	return d.router.onMsgUri(uri, cb)
}

func (d *Dealer) OnReq(uri string, cb func(*types.Request) bool) (unsubscribe func()) {
	return d.router.onReqUri(uri, cb)
}

func (d *Dealer) Start(ctx context.Context) error {
	if !d.running.CompareAndSwap(false, true) {
		return errors.New("dealer: already started")
	}

	ctx, cancel := context.WithCancel(ctx)
	d.cancelLoop = cancel
	d.done = make(chan struct{})

	go d.loop(ctx)

	Subscribe(d, TopicConnectionID, func(s string) {
		d.setConnectionID(s)
	})

	return nil
}

func (d *Dealer) Stop() error {
	return d.Goodbye(context.Background())
}

// Goodbye closes the live websocket gracefully and stops the reconnect loop.
func (d *Dealer) Goodbye(ctx context.Context) error {
	if !d.running.CompareAndSwap(true, false) {
		return nil
	}

	d.connMu.Lock()
	if c := d.conn; c != nil {
		c.goodbye()
	}
	d.connMu.Unlock()

	if d.cancelLoop != nil {
		d.cancelLoop()
	}

	if d.done == nil {
		return nil
	}
	select {
	case <-d.done:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (d *Dealer) Send(msg []byte) error {
	if d.conn == nil {
		return ErrNotConnected
	}

	d.connMu.RLock()
	ch := d.conn.send
	d.connMu.RUnlock()

	select {
	case ch <- msg:
		return nil
	default:
		return ErrSendOverflow
	}
}
