package pathfinder

import (
	"bytes"
	"context"
	"io"
	"net/http"

	"github.com/goccy/go-json"
	"github.com/pkg/errors"
	"github.com/pyrorhythm/fn/bjs"
	"github.com/pyrorhythm/libspot"
	"github.com/pyrorhythm/libspot/auth/client"
	"github.com/pyrorhythm/libspot/pathfinder/responsetypes"
	"github.com/pyrorhythm/libspot/pathfinder/types"
	"github.com/pyrorhythm/libspot/pkg/transport"
)

const (
	reqUrl = "https://api-partner.spotify.com/pathfinder/v2/query"
)

type Pathfinder struct {
	prov   libspot.TokenProvider
	client *http.Client
}

func New(prov libspot.TokenProvider) *Pathfinder {
	return &Pathfinder{
		prov: prov,
		client: client.NewAuthorizedClient(prov).
			BaseTransport(&transport.LoggingTransport{}).
			Client(),
	}
}

type SearchResponse types.RespPayload[responsetypes.SearchResultV2]

func (p *Pathfinder) QuerySuggestions(
	ctx context.Context,
	payload types.SuggestionsPayload,
) (*SearchResponse, error) {
	op := types.OpSuggestions

	bs, err := json.Marshal(&types.ReqPayload[types.SuggestionsPayload]{
		Variables:     payload,
		OperationName: op,
		Extensions: &types.Extensions{
			PersistedQuery: op.Extension(),
		},
	})
	if err != nil {
		return nil, errors.Wrap(err, "failed to marshal payload")
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, reqUrl, bytes.NewReader(bs))
	if err != nil {
		return nil, errors.Wrap(err, "failed to create request")
	}

	resp, err := p.makeRequest(ctx, req)
	if err != nil {
		return nil, err
	}

	defer func() {
		_ = resp.Body.Close()
	}()
	raw, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return bjs.Unmarshal[SearchResponse](raw)
}
