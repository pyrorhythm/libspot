package dealer

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"github.com/pyrorhythm/libspot"
	"github.com/pyrorhythm/libspot/auth/session"
	"github.com/pyrorhythm/libspot/dealer/types"
	"github.com/pyrorhythm/libspot/pkg/delay"
)

var (
	ErrNotConnected = errors.New("dealer: not connected")
	ErrSendOverflow = errors.New("dealer: send buffer full")
)

type DelayConfig struct {
	Fn  delay.Delay
	Cap int64
}

//goland:noinspection GoNameStartsWithPackageName
type Dealer struct {
	prov  libspot.TokenProvider
	rslv  libspot.EndpointResolver
	delay func(attempt int64) time.Duration

	endpoint *DelayConfig
	global   *DelayConfig

	intervalSec time.Duration
	timeout     time.Duration

	onConnectionShutdown func(error)

	router Router

	conn   *conn
	connMu sync.RWMutex

	connectionId string

	running    atomic.Bool
	cancelLoop context.CancelFunc
}

var commonDelayCfg = &DelayConfig{
	Fn:  delay.ExponentialJitter2Delay(2 * time.Second),
	Cap: 5,
}

func applyDefaults(dealer *Dealer) {
	dealer.endpoint = commonDelayCfg
	dealer.global = commonDelayCfg

	dealer.intervalSec = 10 * time.Second
	dealer.timeout = 30 * time.Second

	dealer.onConnectionShutdown = func(error) {}
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

	if dealer.intervalSec <= 0 {
		dealer.intervalSec = 10 * time.Second
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

	applyOptions(d, opts)

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

	d := &Dealer{
		prov:   sess,
		rslv:   rslv,
		router: newRouter(),
	}

	applyOptions(d, opts)

	return d, nil
}

func applyOptions(d *Dealer, opts []Option) {
	applyDefaults(d)

	for _, opt := range opts {
		opt(d)
	}

	coverNils(d)
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

	go d.loop(ctx)

	Subscribe(d, TopicConnectionID, func(s string) {
		d.connectionId = s
	})

	return nil
}

func (d *Dealer) Stop() error {
	if !d.running.CompareAndSwap(true, false) {
		return nil
	}
	d.cancelLoop()
	return nil
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
