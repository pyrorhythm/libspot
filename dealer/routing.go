package dealer

import (
	"context"
	"strings"
	"sync"
	"sync/atomic"
)

type router struct {
	mrMu sync.RWMutex
	mr   map[string]*msgRouter

	reqMu sync.RWMutex
	req   map[string]func(*Msg) bool
}

type msgRouter struct {
	id   atomic.Uint64
	mu   sync.RWMutex
	subs map[uint64]func(*Msg)
}

type unsubFn func()

func (r *msgRouter) notify(msg *Msg) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	for _, v := range r.subs {
		v(msg)
	}
}

func (r *msgRouter) add(cb func(*Msg)) unsubFn {
	r.mu.Lock()
	defer r.mu.Unlock()
	idx := r.id.Add(1)
	r.subs[idx] = cb
	return func() {
		r.mu.Lock()
		defer r.mu.Unlock()
		delete(r.subs, idx)
	}
}

type Router interface {
	handleMsg(msg *Msg)
	handleReq(ctx context.Context, msg *Msg) (reply chan bool, found bool)
	onMsgUri(uri string, cb func(*Msg)) unsubFn
	onReqUri(uri string, cb func(*Msg) bool) unsubFn
}

func newRouter() Router {
	return &router{
		mr:  make(map[string]*msgRouter),
		req: make(map[string]func(*Msg) bool),
	}
}

func (r *router) handleMsg(msg *Msg) {
	r.mrMu.RLock()
	defer r.mrMu.RUnlock()
	for uri, mr := range r.mr {
		if strings.HasPrefix(msg.Uri, uri) {
			// TODO make better
			go mr.notify(msg)
		}
	}
}

func (r *router) handleReq(ctx context.Context, msg *Msg) (replyChan chan bool, found bool) {
	var reqsub struct {
		plen int
		cb   func(*Msg) bool
	}

	r.reqMu.RLock()
	for uri, reqCb := range r.req {
		if strings.HasPrefix(msg.Uri, uri) && len(uri) > reqsub.plen {
			reqsub.plen = len(uri)
			reqsub.cb = reqCb
		}
	}
	r.reqMu.RUnlock()

	if reqsub.cb == nil {
		return nil, false
	}

	replyChan = make(chan bool, 1)

	go func() {
		defer close(replyChan)
		select {
		case <-ctx.Done():
			return
		case replyChan <- reqsub.cb(msg):
			return
		}
	}()

	return replyChan, true
}

func (r *router) onReqUri(uri string, cb func(*Msg) bool) unsubFn {
	r.reqMu.Lock()
	defer r.reqMu.Unlock()

	if _, ok := r.req[uri]; ok {
		panic("libspot: dealer request can have only 1 sub")
	}

	r.req[uri] = cb

	return func() {
		r.reqMu.Lock()
		defer r.reqMu.Unlock()

		delete(r.req, uri)
	}
}

func (r *router) onMsgUri(uri string, cb func(*Msg)) unsubFn {
	r.mrMu.Lock()
	defer r.mrMu.Unlock()

	mr, ok := r.mr[uri]
	if !ok {
		mr = &msgRouter{subs: make(map[uint64]func(*Msg))}
	}

	return mr.add(cb)
}
