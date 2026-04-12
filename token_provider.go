package libspot

import "context"

type TokenProvider interface {
	ClientToken() (string, error)
	AccessToken(ctx context.Context, refresh bool) (string, error)
}
