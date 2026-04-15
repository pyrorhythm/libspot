package pathfinder

import (
	"github.com/pyrorhythm/libspot"
	"github.com/pyrorhythm/libspot/auth/client"
	"github.com/pyrorhythm/libspot/pkg/transport"
	"resty.dev/v3"
)

const (
	reqUrl = "https://api-partner.spotify.com/pathfinder/v2/query"
)

type Pathfinder struct {
	prov   libspot.TokenProvider
	client *resty.Client
}

func New(prov libspot.TokenProvider) *Pathfinder {
	return &Pathfinder{
		prov: prov,
		client: client.NewAuthorizedClient(prov).
			BaseTransport(&transport.LoggingTransport{}).
			Client(),
	}
}
