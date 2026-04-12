package session

import (
	"crypto/ecdsa"
	"net/http"

	"github.com/pyrorhythm/libspot/pkg/dpop"
	"github.com/pyrorhythm/libspot/pkg/transport"
)

func newSessionDpopClient(
	clientToken string,
	sess *session,
	key *ecdsa.PrivateKey,
) *http.Client {
	return &http.Client{
		Transport: &clientTokenTransport{
			clientToken: clientToken,
			base: &dpop.Transport{
				Base: &transport.LoggingTransport{},
				Prov: &dpop.Provider{
					GetAccessToken: sess.GetToken,
					GetNonce:       sess.GetNonce,
					SetNonce:       sess.SetNonce,
				},
				Key: key,
			},
		},
	}
}

type clientTokenTransport struct {
	clientToken string
	base        http.RoundTripper
}

func (c *clientTokenTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	if c.clientToken != "" {
		req = req.Clone(req.Context())
		req.Header.Set("Client-Token", c.clientToken)
	}
	return c.base.RoundTrip(req)
}
