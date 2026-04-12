// Package spcli provides Spclient implementations.
// Spotify has two types of them right now: spclient.wg.(spotify.com),
// and g(region)(id)-spclient.(...) (gew4-spclient as an example).
// First one is doing mostly JSON-oriented requests, second one provides some
// key protobuf endpoints. Both are implemented in a client below, as separate services
package spcli

type Spclient interface {
	stub()
}
