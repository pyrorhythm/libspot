package client

import (
	"net/http"

	"resty.dev/v3"
)

type authorizedClientBuilder interface {
	InjectClientToken(b bool) authorizedClientBuilder
	InjectAccessToken(b bool) authorizedClientBuilder
	CanRefreshAccessToken(b bool) authorizedClientBuilder
	BaseTransport(http.RoundTripper) authorizedClientBuilder

	Transport() http.RoundTripper
	Client() *resty.Client
}
