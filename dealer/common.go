package dealer

import (
	"context"
	"errors"
	"fmt"
	"math"
	"math/rand"
	"sync"
	"sync/atomic"
	"time"

	"github.com/pyrorhythm/libspot"
	"github.com/pyrorhythm/libspot/auth/session"
	dtyp "github.com/pyrorhythm/libspot/dealer/types"
)

var (
	ErrNotConnected = errors.New("dealer: not connected")
	ErrSendOverflow = errors.New("dealer: send buffer full")
)

func LinearDelay(base time.Duration, increment time.Duration) func(int64) time.Duration {
	return func(att int64) time.Duration {
		return base + increment*time.Duration(att-1)
	}
}

// ExponentialDelay / powBase -- seconds
func ExponentialDelay(base time.Duration, powBase int64) func(int64) time.Duration {
	return func(att int64) time.Duration {
		return base + float64DurSec(math.Pow(float64(powBase), float64(att-1))) * time.Second
	}
}

func Exponential2Delay(base time.Duration) func(int64) time.Duration {
	return func(att int64) time.Duration {
		return base + float64DurSec(math.Pow(2, float64(att-1)))
	}
}

func ExponentialJitterDelay(base time.Duration, pow int64) func(int64) time.Duration {
	return func(att int64) time.Duration {
		return base + float64DurSec(math.Pow(float64(pow), float64(att-1))*(0.5+0.5*rand.Float64()))
	}
}

func ExponentialJitter2Delay(base time.Duration) func(int64) time.Duration {
	return func(att int64) time.Duration {
		return base + float64DurSec(math.Pow(2, float64(att-1))*(0.5+0.5*rand.Float64()))
	}
}

func float64DurSec(f float64) time.Duration {
	return time.Duration(f) * time.Second
}

//goland:noinspection GoNameStartsWithPackageName
type Dealer struct {
	prov  libspot.TokenProvider
	rslv  libspot.EndpointResolver
	delay func(attempt int64) time.Duration

	RetryCap         int64
	GlobalRetryDelay int64 // seconds
	PingIntervalSec  int64
	PingTimeout      int64

	OnClose func(err error)

	router Router // main Router instance, persisted through connections

	conn   *conn
	connMu sync.RWMutex
	
	connectionId string

	running atomic.Bool
	loopCl  context.CancelFunc
}

func New(
	prov libspot.TokenProvider,
	rslv libspot.EndpointResolver,
	delay func(attempt int64) time.Duration,
) *Dealer {
	return &Dealer{
		prov:   prov,
		rslv:   rslv,
		delay:  delay,
		router: newRouter(),
	}
}

func NewSession(
	sess session.Session,
	delay func(attempt int64) time.Duration,
) (*Dealer, error) {
	rslv, err := sess.Resolver()
	if err != nil {
		return nil, fmt.Errorf("failed to get resolver from session: %w", err)
	}

	return &Dealer{
		prov:   sess,
		rslv:   rslv,
		delay:  delay,
		router: newRouter(),
	}, nil
}

func (d *Dealer) OnMsg(uri string, cb func(*dtyp.Message)) (unsubscribe func()) {
	return d.router.onMsgUri(uri, cb)
}

func (d *Dealer) OnReq(uri string, cb func(*dtyp.Request) bool) (unsubscribe func()) {
	return d.router.onReqUri(uri, cb)
}

func (d *Dealer) Start() error {
	if !d.running.CompareAndSwap(false, true) {
		return errors.New("dealer: already started")
	}

	ctx, cancel := context.WithCancel(context.Background())
	d.loopCl = cancel

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
	d.loopCl()
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
