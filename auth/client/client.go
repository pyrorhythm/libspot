package client

import (
	"net/http"

	"pyrorhythm.dev/libspot"
	libtransport "pyrorhythm.dev/libspot/pkg/transport"
	"resty.dev/v3"
)

type transportSettings struct {
	injectClientToken     bool
	injectAccessToken     bool
	canRefreshAccessToken bool
	baseTransport         http.RoundTripper
}

type AuthorizedClient struct {
	prov         libspot.TokenProvider
	baseSettings *transportSettings
}

func defaults(ac *AuthorizedClient) {
	ac.baseSettings = &transportSettings{
		injectClientToken:     true,
		injectAccessToken:     true,
		canRefreshAccessToken: false,
	}
}

func applyOptions[T any, O func(*T)](p *T, defaults O, opts ...O) *T {
	defaults(p)

	for _, opt := range opts {
		opt(p)
	}

	return p
}

func NewAuthorizedClient(
	prov libspot.TokenProvider,
	opts ...func(*AuthorizedClient),
) *AuthorizedClient {
	return applyOptions(&AuthorizedClient{prov: prov}, defaults, opts...)
}

func (a *AuthorizedClient) Transport() http.RoundTripper {
	setts := *a.baseSettings

	if setts.baseTransport == nil {
		setts.baseTransport = libtransport.RoundTripper()
	}

	return &transport{prov: a.prov, transportSettings: setts}
}

func (a *AuthorizedClient) Client() *resty.Client {
	rc := resty.New()
	rc.SetTransport(a.Transport())
	return rc
}

func (a *AuthorizedClient) InjectClientToken(b bool) authorizedClientBuilder {
	bd := &authorizedCBImpl{prov: a.prov, setts: *a.baseSettings}
	bd.setts.injectClientToken = b
	return bd
}

func (a *AuthorizedClient) InjectAccessToken(b bool) authorizedClientBuilder {
	bd := &authorizedCBImpl{prov: a.prov, setts: *a.baseSettings}
	bd.setts.injectAccessToken = b
	return bd
}

func (a *AuthorizedClient) CanRefreshAccessToken(b bool) authorizedClientBuilder {
	bd := &authorizedCBImpl{prov: a.prov, setts: *a.baseSettings}
	bd.setts.canRefreshAccessToken = b
	return bd
}

func (a *AuthorizedClient) BaseTransport(rt http.RoundTripper) authorizedClientBuilder {
	bd := &authorizedCBImpl{prov: a.prov, setts: *a.baseSettings}
	bd.setts.baseTransport = rt
	return bd
}
