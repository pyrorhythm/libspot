package libspot

import "context"

type TokenProvider interface {
	ClientToken() (string, error)
	GetOrRefreshToken(ctx context.Context) (string, error)
}
