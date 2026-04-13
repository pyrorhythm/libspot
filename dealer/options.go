package dealer

import "time"

type Option func(*Dealer)

func WithEndpointDelayConfig(cfg *DelayConfig) Option {
	return func(dealer *Dealer) {
		dealer.endpoint = cfg
	}
}

func WithGlobalDelayConfig(cfg *DelayConfig) Option {
	return func(dealer *Dealer) {
		dealer.global = cfg
	}
}

func WithPingInterval(interval time.Duration) Option {
	return func(dealer *Dealer) {
		dealer.intervalSec = interval
	}
}

func WithPingTimeout(timeout time.Duration) Option {
	return func(dealer *Dealer) {
		dealer.timeout = timeout
	}
}

func WithOnConnectionShutdown(f func(error)) Option {
	return func(dealer *Dealer) {
		dealer.onConnectionShutdown = f
	}
}
