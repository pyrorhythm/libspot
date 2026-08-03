package client

import (
	"net/http"

	"pyrorhythm.dev/libspot"
	libtransport "pyrorhythm.dev/libspot/pkg/transport"
	"resty.dev/v3"
)

type authorizedCBImpl struct {
	prov  libspot.TokenProvider
	setts transportSettings
}

func (a *authorizedCBImpl) BaseTransport(rt http.RoundTripper) authorizedClientBuilder {
	a.setts.baseTransport = rt
	return a
}

func (a *authorizedCBImpl) InjectClientToken(b bool) authorizedClientBuilder {
	a.setts.injectClientToken = b
	return a
}

func (a *authorizedCBImpl) InjectAccessToken(b bool) authorizedClientBuilder {
	a.setts.injectAccessToken = b
	return a
}

func (a *authorizedCBImpl) CanRefreshAccessToken(b bool) authorizedClientBuilder {
	a.setts.canRefreshAccessToken = b
	return a
}

func (a *authorizedCBImpl) Transport() http.RoundTripper {
	if a.setts.baseTransport == nil {
		a.setts.baseTransport = libtransport.RoundTripper()
	}

	return &transport{prov: a.prov, transportSettings: a.setts}
}

func (a *authorizedCBImpl) Client() *resty.Client {
	rc := resty.New()
	rc.SetTransport(a.Transport())
	return rc
}
