// Package pbsp provides Protocol Buffers based spclient
package pbsp

import (
	"context"

	"github.com/pyrorhythm/libspot"
)

type PBClient struct {
	prov    libspot.TokenProvider
	resolve libspot.EndpointResolver
}

func New(
	prov libspot.TokenProvider,
	resolve libspot.EndpointResolver,
) *PBClient {
	return &PBClient{
		prov:    prov,
		resolve: resolve,
	}
}

type SpclientProto interface {
	// extended-metadata/v0
	ExtendedMetadata(ctx context.Context)
	// connect-state/v1
	//
}
