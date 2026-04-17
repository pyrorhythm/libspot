package pathfinder

import (
	"context"
	"net/http"

	"github.com/pkg/errors"
	pfq "github.com/pyrorhythm/libspot/pathfinder/pfrequest"
	pfs "github.com/pyrorhythm/libspot/pathfinder/pfresponse"
)

func (p *Pathfinder) searchV2(ctx context.Context, rq pfq.Request) (*pfs.Response, error) {
	bs, err := Marshal(rq)
	if err != nil {
		return nil, errors.Wrap(err, "failed to marshal payload")
	}

	resp, err := p.makeRequest(ctx, http.MethodPost, nil, bs)
	if err != nil {
		return nil, err
	}

	sresp, err := Unmarshal[pfs.Response](resp.Bytes())
	if err != nil {
		return nil, err
	}

	return sresp.Get(), nil
}

func (p *Pathfinder) Query(
	ctx context.Context,
	rq pfq.Request,
) (*pfs.Response, error) {
	return p.searchV2(ctx, rq)
}
