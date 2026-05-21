package spc

import (
	"context"
	"log/slog"
	"net/http"
	"net/url"

	"github.com/pkg/errors"
	"github.com/pyrorhythm/libspot"
	extmetadatapb "github.com/pyrorhythm/libspot/gen/spotify/extendedmetadata"
	"github.com/pyrorhythm/libspot/spc/extmetadata"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

// ExtendedMetadata fetches batched per-entity extensions (audio files, canvas,
// track descriptors, ...) from the extended-metadata endpoint. Build req with the
// extmetadata package to choose which entities and ExtensionKinds to request:
//
//	resp, err := c.ExtendedMetadata(ctx, extmetadata.New().
//		Country("US").Catalogue("premium").
//		Query("spotify:track:"+id, extmetadata.KindAudioFiles, extmetadata.KindCanvaz))
func (c *Spclient) ExtendedMetadata(
	ctx context.Context,
	req *extmetadata.Request,
) (*extmetadatapb.BatchedExtensionResponse, error) {
	body, err := proto.Marshal(req.Build())
	if err != nil {
		return nil, errors.Wrap(err, "failed to marshal extended metadata request")
	}

	u, _ := url.JoinPath(baseUrl, extmetadata.Path)

	resp := &extmetadatapb.BatchedExtensionResponse{}
	if err := makeProtoRequest(ctx,
		c.client.R().SetURL(u).SetMethod(http.MethodPost).SetBody(body),
		resp,
	); err != nil {
		return nil, errors.Wrap(err, "failed to fetch extended metadata")
	}

	slog.Debug("extended metadata", "resp", protojson.Format(resp))

	return resp, nil
}

// ExtendedMetadataSimple is a convenience over ExtendedMetadata for the common
// case of fetching a single extension kind for one entity. It locates the
// matching extension in the batched response and unmarshals its payload into
// data, which must be the concrete protobuf type for ext.
func (c *Spclient) ExtendedMetadataSimple(
	ctx context.Context,
	id libspot.SpotifyId,
	ext extmetadata.ExtensionKind,
	tgt proto.Message,
) error {
	resp, err := c.ExtendedMetadata(ctx, extmetadata.New().Query(id.Uri(), ext))
	if err != nil {
		return err
	}

	for _, item := range resp.GetExtendedMetadata() {
		if item.GetExtensionKind() != ext {
			continue
		}

		for _, extData := range item.GetExtensionData() {
			if extData.GetEntityUri() != id.Uri() {
				continue
			}

			if code := extData.GetHeader().GetStatusCode(); code != http.StatusOK {
				return errors.Errorf("extended metadata request returned status %d", code)
			}

			if err := extData.GetExtensionData().UnmarshalTo(tgt); err != nil {
				return errors.Wrap(err, "failed to unmarshal extended metadata data")
			}

			return nil
		}
	}

	return errors.Errorf("extended metadata with kind %s not found", ext)
}
