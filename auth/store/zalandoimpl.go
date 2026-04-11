package store

import (
	"encoding/base64"
	"errors"
	"fmt"

	"github.com/bytedance/sonic"
	"github.com/pyrorhythm/fn/errs"
	"github.com/zalando/go-keyring"
)

var _ Keychainer[any] = (*zalandoKeychainer[any])(nil)

type zalandoKeychainer[T any] struct {
	key string

	cached *T
}

func NewZalandoKeychainer[T any](key string) Keychainer[T] {
	return &zalandoKeychainer[T]{key: key}
}

func (z zalandoKeychainer[T]) Load(invalidate bool) (*T, error) {
	if !invalidate && z.cached != nil {
		return z.cached, nil
	}

	data, err := keyring.Get(storeService, z.key)
	if err != nil && !errors.Is(err, keyring.ErrNotFound) {
		return nil, fmt.Errorf("failed to get: %w", err)
	} else if errors.Is(err, keyring.ErrNotFound) {
		return nil, ErrItemNotFound
	}

	bytes, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		return nil, fmt.Errorf("failed to base64 decode: %w", err)
	}

	z.cached = new(T)
	if err := sonic.Unmarshal(bytes, z.cached); err != nil {
		return nil, fmt.Errorf("failed to unmarshal: %w", err)
	}

	return z.cached, nil
}

func (z zalandoKeychainer[T]) Save(item *T) error {
	bytes, err := sonic.Marshal(item)
	if err != nil {
		return fmt.Errorf("failed to marshal: %w", err)
	}

	data := base64.StdEncoding.EncodeToString(bytes)

	return errs.Wrap(keyring.Set(storeService, z.key, data), "failed to set")
}

func (z zalandoKeychainer[T]) Delete() error {
	if err := keyring.Delete(storeService, z.key); err != nil {
		return fmt.Errorf("failed to delete: %w", err)
	}
	z.cached = nil
	return nil
}
