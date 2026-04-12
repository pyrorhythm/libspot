package dealer

import (
	"context"
	"sync"
	"sync/atomic"
	"time"

	"github.com/pyrorhythm/libspot"
)

type ConnState int

const (
	ConnStateIdle ConnState = iota
	ConnStateConnecting
	ConnStateConnected
	ConnStateDisconnecting
)

//goland:noinspection GoNameStartsWithPackageName
type DealerG2 struct {
	prov  libspot.TokenProvider
	rslv  libspot.EndpointResolver
	delay func(startDelay time.Duration, attempt uint) time.Duration

	RetryCap         uint
	GlobalRetryDelay int // seconds
	PingIntervalSec  int
	PingTimeout      int

	OnClose func(err error)

	router Router

	mu     sync.RWMutex
	state  ConnState
	sendCh chan<- []byte

	running       atomic.Bool
	runningCtx    context.Context
	runningCancel context.CancelFunc
}

func NewDealerG2(
	prov libspot.TokenProvider,
	rslv libspot.EndpointResolver,
	delay func(startDelay time.Duration, attempt uint) time.Duration,
) *DealerG2 {
	return &DealerG2{
		prov:   prov,
		rslv:   rslv,
		delay:  delay,
		router: newRouter(),
	}
}

func (d *DealerG2) OnMsg(uri string, cb func(*Msg)) (unsubscribe func()) {
	return d.router.onMsgUri(uri, cb)
}

func (d *DealerG2) OnReq(uri string, cb func(*Msg) bool) (unsubscribe func()) {
	return d.router.onReqUri(uri, cb)
}

func (d *DealerG2) State() ConnState {
	d.mu.RLock()
	defer d.mu.RUnlock()
	return d.state
}
