package client

import "net/http"

func InjectClientToken(b bool) func(*AuthorizedClient) {
	return func(ac *AuthorizedClient) {
		ac.baseSettings.injectClientToken = b
	}
}

func InjectAccessToken(b bool) func(*AuthorizedClient) {
	return func(ac *AuthorizedClient) {
		ac.baseSettings.injectAccessToken = b
	}
}

func CanRefreshAccessToken(b bool) func(*AuthorizedClient) {
	return func(ac *AuthorizedClient) {
		ac.baseSettings.canRefreshAccessToken = b
	}
}

func BaseTransport(tpt http.RoundTripper) func(*AuthorizedClient) {
	return func(ac *AuthorizedClient) {
		ac.baseSettings.baseTransport = tpt
	}
}
