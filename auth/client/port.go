package client

import (
	"net/http"
)

type authorizedClientBuilder interface {
	InjectClientToken(b bool) authorizedClientBuilder
	InjectAccessToken(b bool) authorizedClientBuilder
	CanRefreshAccessToken(b bool) authorizedClientBuilder
	BaseTransport(http.RoundTripper) authorizedClientBuilder

	Transport() http.RoundTripper
	Client() *http.Client
}
