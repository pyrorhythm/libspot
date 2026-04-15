package pathfinder

import (
	"context"
	"net/http"

	"github.com/goccy/go-json"
	"github.com/pkg/errors"
	"github.com/pyrorhythm/fn/bjs"
	pfd "github.com/pyrorhythm/libspot/pathfinder/pfdomain"
	pfreq "github.com/pyrorhythm/libspot/pathfinder/pfrequest"
	pfresp "github.com/pyrorhythm/libspot/pathfinder/pfresponse"
)

func (p *Pathfinder) QuerySuggestions(
	ctx context.Context,
	payload pfreq.SuggestionsPayload,
) (*pfresp.SearchResultV2, error) {
	op := pfreq.OpSuggestions

	bs, err := json.Marshal(&pfreq.Payload[pfreq.SuggestionsPayload]{
		Variables:     payload,
		OperationName: op,
		Extensions: &pfd.Extensions{
			PersistedQuery: op.Extension(),
		},
	})
	if err != nil {
		return nil, errors.Wrap(err, "failed to marshal payload")
	}

	resp, err := p.makeRequest(ctx, http.MethodPost, nil, bs)
	if err != nil {
		return nil, err
	}

	sresp, err := bjs.Unmarshal[pfresp.Payload[pfresp.SearchResultV2]](resp.Bytes())
	if err != nil {
		return nil, err
	}

	return sresp.Data, nil
}
