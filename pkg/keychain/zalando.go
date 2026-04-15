package keychain

import (
	"encoding/base64"
	"errors"
	"fmt"

	"github.com/goccy/go-json"
	"github.com/pyrorhythm/fn/errs"
	"github.com/zalando/go-keyring"
)

var _ Keychainer[any] = (*zalandoKeychainer[any])(nil)

type zalandoKeychainer[T any] struct {
	service string
	key     string
	cached  *T
}

func NewZalandoKeychainerSvcFn[T any](service string) func(string) Keychainer[T] {
	return func(key string) Keychainer[T] { return &zalandoKeychainer[T]{service: service, key: key} }
}

func NewZalandoKeychainer[T any](service, key string) Keychainer[T] {
	return &zalandoKeychainer[T]{service: service, key: key}
}

func (z *zalandoKeychainer[T]) Load(invalidate bool) (*T, error) {
	if !invalidate && z.cached != nil {
		return z.cached, nil
	}

	data, err := keyring.Get(z.service, z.key)
	if err != nil && !errors.Is(err, keyring.ErrNotFound) {
		return nil, fmt.Errorf("failed to get: %w", err)
	} else if errors.Is(err, keyring.ErrNotFound) {
		return nil, ErrItemNotFound
	}

	b, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		return nil, fmt.Errorf("failed to base64 decode: %w", err)
	}

	z.cached = new(T)
	if err := json.Unmarshal(b, z.cached); err != nil {
		return nil, fmt.Errorf("failed to unmarshal: %w", err)
	}
	return z.cached, nil
}

func (z *zalandoKeychainer[T]) Save(item *T) error {
	b, err := json.Marshal(item)
	if err != nil {
		return fmt.Errorf("failed to marshal: %w", err)
	}
	data := base64.StdEncoding.EncodeToString(b)
	return errs.Wrap(keyring.Set(z.service, z.key, data), "failed to set")
}

func (z *zalandoKeychainer[T]) Delete() error {
	if err := keyring.Delete(z.service, z.key); err != nil {
		return fmt.Errorf("failed to delete: %w", err)
	}
	z.cached = nil
	return nil
}
