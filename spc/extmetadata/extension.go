package extmetadata

import (
	"context"
	"net/http"

	"github.com/pkg/errors"
	"github.com/pyrorhythm/libspot"
	pb "github.com/pyrorhythm/libspot/gen/spotify/extendedmetadata"
	"google.golang.org/protobuf/proto"
)

type Client interface {
	ExtendedMetadata(ctx context.Context, req *Request) (*pb.BatchedExtensionResponse, error)
}

type protoPtr[T any] interface {
	*T
	proto.Message
}

// Extension declaratively binds an ExtensionKind to the concrete protobuf
// message its payload unmarshals into. Declare one with Define (or use the
// predefined vars below) and call Fetch to retrieve and decode it in one step:
//
//	files, err := extmetadata.AudioFiles.Fetch(ctx, client, id)
type Extension[T any, PT protoPtr[T]] struct {
	kind ExtensionKind
}

// Define declares a typed extension descriptor. Specify the message value type;
// the pointer type is inferred:
//
//	var AudioFiles = Define[audiofiles.AudioFilesExtensionResponse](KindAudioFiles)
func Define[T any, PT protoPtr[T]](kind ExtensionKind) Extension[T, PT] {
	return Extension[T, PT]{kind: kind}
}

// Fetch requests this extension for a single entity and unmarshals the matching
// payload. It errors if the entity is missing from the response, the per-entity
// status is non-200, or the payload type mismatches.
func (e Extension[T, PT]) Fetch(ctx context.Context, c Client, id libspot.SpotifyId) (PT, error) {
	resp, err := c.ExtendedMetadata(ctx, New().Query(id.Uri(), e.kind))
	if err != nil {
		return nil, err
	}

	uri := id.Uri()
	for _, item := range resp.GetExtendedMetadata() {
		if item.GetExtensionKind() != e.kind {
			continue
		}

		for _, data := range item.GetExtensionData() {
			if data.GetEntityUri() != uri {
				continue
			}

			if code := data.GetHeader().GetStatusCode(); code != http.StatusOK {
				return nil, errors.Errorf("extended metadata %s returned status %d", e.kind, code)
			}

			out := PT(new(T))
			if err := data.GetExtensionData().UnmarshalTo(out); err != nil {
				return nil, errors.Wrapf(err, "failed to unmarshal %s payload", e.kind)
			}

			return out, nil
		}
	}

	return nil, errors.Errorf("extended metadata %s not found for %s", e.kind, uri)
}
