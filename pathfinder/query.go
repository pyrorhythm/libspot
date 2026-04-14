package pathfinder

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/bytedance/sonic"
	"github.com/pyrorhythm/libspot"
	"github.com/pyrorhythm/libspot/auth/client"
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
		client: client.NewAuthorizedClient(
			prov, client.BaseTransport(&transport.LoggingTransport{})).
			Client(),
	}
}

func (p *Pathfinder) QuerySuggestions(
	ctx context.Context,
	payload types.SuggestionsPayload,
) ([]byte, error) {
	rurl, err := url.Parse(reqUrl)
	if err != nil {
		return nil, fmt.Errorf("failed to parse url: %w", err)
	}

	bs, err := sonic.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal payload: %w", err)
	}

	req := (&http.Request{
		URL:        rurl,
		Method:     http.MethodPost,
		Body:       io.NopCloser(bytes.NewReader(bs)),
		Proto:      "HTTP/2.0",
		ProtoMajor: 2,
		ProtoMinor: 0,
		Header:     make(http.Header),
		Host:       rurl.Host,
	}).WithContext(ctx)

	resp, err := p.makeRequest(ctx, req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	return io.ReadAll(resp.Body)
}
