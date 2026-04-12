package libspot

import "context"

type TokenProvider interface {
	AccessToken(ctx context.Context) (string, error)
}
