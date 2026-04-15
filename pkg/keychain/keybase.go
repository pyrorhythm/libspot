package keychain

import (
	"errors"
	"fmt"

	"github.com/goccy/go-json"
	"github.com/keybase/go-keychain"
)

var _ Keychainer[any] = (*keybaseKeychainer[any])(nil)

type keybaseKeychainer[T any] struct {
	service string
	key     string
	cached  *T
}

func NewKeybaseKeychainer[T any](service, key string) Keychainer[T] {
	return &keybaseKeychainer[T]{service: service, key: key}
}

func (k *keybaseKeychainer[T]) Load(invalidate bool) (*T, error) {
	if !invalidate && k.cached != nil {
		return k.cached, nil
	}

	data, err := keychain.GetGenericPassword(k.service, k.key, "", "")
	if err != nil || data == nil {
		if errors.Is(err, keychain.ErrorItemNotFound) || data == nil {
			return nil, ErrItemNotFound
		}
		return nil, fmt.Errorf("keychain query failed: %w", err)
	}

	k.cached = new(T)
	if err := json.Unmarshal(data, k.cached); err != nil {
		return nil, fmt.Errorf("failed to unmarshal object %T: %w", *new(T), err)
	}
	return k.cached, nil
}

func (k *keybaseKeychainer[T]) Save(item *T) error {
	data, err := json.Marshal(item)
	if err != nil {
		return fmt.Errorf("failed to marshal: %w", err)
	}

	_, err = keychain.GetGenericPassword(k.service, k.key, "", "")
	if err != nil && !errors.Is(err, keychain.ErrorItemNotFound) {
		return fmt.Errorf("failed to get generic password: %w", err)
	}

	if err == nil {
		return k.updateItem(data)
	}
	return k.addItem(data)
}

func (k *keybaseKeychainer[T]) addItem(data []byte) error {
	item := keychain.NewGenericPassword(k.service, k.key, "", data, "")
	item.SetAccessible(keychain.AccessibleWhenUnlockedThisDeviceOnly)
	if err := keychain.AddItem(item); err != nil {
		return fmt.Errorf("keychain add failed: %w", err)
	}
	return nil
}

func (k *keybaseKeychainer[T]) updateItem(data []byte) error {
	updateItem := keychain.NewItem()
	updateItem.SetData(data)
	updateItem.SetAccessible(keychain.AccessibleWhenUnlockedThisDeviceOnly)
	if err := keychain.UpdateItem(k.queryItem(), updateItem); err != nil {
		return fmt.Errorf("keychain update failed: %w", err)
	}
	return nil
}

func (k *keybaseKeychainer[T]) queryItem() keychain.Item {
	query := keychain.NewItem()
	query.SetSecClass(keychain.SecClassGenericPassword)
	query.SetService(k.service)
	query.SetAccount(k.key)
	query.SetMatchLimit(keychain.MatchLimitOne)
	query.SetReturnData(true)
	return query
}

func (k *keybaseKeychainer[T]) Delete() error {
	err := keychain.DeleteGenericPasswordItem(k.service, k.key)
	if err != nil && !errors.Is(err, keychain.ErrorItemNotFound) {
		return fmt.Errorf("keychain delete failed: %w", err)
	}
	k.cached = nil
	return nil
}
