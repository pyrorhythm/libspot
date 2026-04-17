// Package pbsp provides Protocol Buffers based spclient
package pbsp

import (
	"context"

	"github.com/pyrorhythm/libspot"
	emd "github.com/pyrorhythm/libspot/gen/spotify/extendedmetadata"
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
	ExtendedMetadata(ctx context.Context, req *emd.BatchedEntityRequest) (*emd.BatchedExtensionResponse, error)
	// connect-state/v1
	//
}
