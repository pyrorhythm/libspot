package store

import (
	"errors"
	"fmt"

	"github.com/bytedance/sonic"
	"github.com/keybase/go-keychain"
)

type keybaseKeychainer[T any] struct {
	key string

	cached *T
}

func NewKeybaseKeychainer[T any](key string) Keychainer[T] {
	return &keybaseKeychainer[T]{
		key: key,
	}
}

func (k keybaseKeychainer[T]) Load(invalidate bool) (*T, error) {
	if !invalidate && k.cached != nil {
		return k.cached, nil
	}

	data, err := keychain.GetGenericPassword(storeService, k.key, "", "")
	if err != nil || data == nil {
		if errors.Is(err, keychain.ErrorItemNotFound) || data == nil {
			return nil, ErrItemNotFound
		}
		return nil, fmt.Errorf("keychain query failed: %w", err)
	}

	k.cached = new(T)
	if err := sonic.Unmarshal(data, k.cached); err != nil {
		return nil, fmt.Errorf("failed to unmarshal object %T: %w", *new(T), err)
	}
	return k.cached, nil
}

func (k keybaseKeychainer[T]) Save(item *T) error {
	data, err := sonic.Marshal(item)
	if err != nil {
		return fmt.Errorf("failed to marshal session: %w", err)
	}

	_, err = keychain.GetGenericPassword(storeService, k.key, "", "")
	if err != nil && !errors.Is(err, keychain.ErrorItemNotFound) {
		return fmt.Errorf("failed to get generic password: %w", err)
	}

	if err == nil {
		return k.updateItem(data)
	}

	return k.addItem(data)
}

func (k keybaseKeychainer[T]) addItem(data []byte) error {
	kcItem := keychain.NewGenericPassword(storeService, k.key, "", data, "")
	kcItem.SetAccessible(keychain.AccessibleWhenUnlockedThisDeviceOnly)
	if aerr := keychain.AddItem(kcItem); aerr != nil {
		return fmt.Errorf("keychain add failed: %w", aerr)
	}
	return nil
}

func (k keybaseKeychainer[T]) updateItem(data []byte) error {
	updateItem := keychain.NewItem()
	updateItem.SetData(data)
	updateItem.SetAccessible(keychain.AccessibleWhenUnlockedThisDeviceOnly)
	if uerr := keychain.UpdateItem(k.queryItem(), updateItem); uerr != nil {
		return fmt.Errorf("keychain update failed: %w", uerr)
	}
	return nil
}

func (k keybaseKeychainer[T]) queryItem() keychain.Item {
	query := keychain.NewItem()
	query.SetSecClass(keychain.SecClassGenericPassword)
	query.SetService(storeService)
	query.SetAccount(k.key)
	query.SetMatchLimit(keychain.MatchLimitOne)
	query.SetReturnData(true)
	return query
}

func (k keybaseKeychainer[T]) Delete() error {
	err := keychain.DeleteGenericPasswordItem(storeService, k.key)
	if err != nil && !errors.Is(err, keychain.ErrorItemNotFound) {
		return fmt.Errorf("keychain delete failed: %w", err)
	}
	k.cached = nil
	return nil
}

func Save[T any](key string, obj T) error {
	data, err := sonic.Marshal(obj)
	if err != nil {
		return fmt.Errorf("failed to marshal session: %w", err)
	}

	query := keychain.NewItem()
	query.SetSecClass(keychain.SecClassGenericPassword)
	query.SetService(storeService)
	query.SetAccount(key)
	query.SetMatchLimit(keychain.MatchLimitOne)
	query.SetReturnData(true)

	_, err = keychain.GetGenericPassword(storeService, key, "", "")
	if err != nil && !errors.Is(err, keychain.ErrorItemNotFound) {
		return fmt.Errorf("failed to get generic password: %w", err)
	} else if err == nil {
		updateItem := keychain.NewItem()
		updateItem.SetData(data)
		updateItem.SetAccessible(keychain.AccessibleWhenUnlockedThisDeviceOnly)
		if uerr := keychain.UpdateItem(query, updateItem); uerr != nil {
			return fmt.Errorf("keychain update failed: %w", uerr)
		}
		return nil
	}

	// Add new item
	item := keychain.NewGenericPassword(storeService, key, "", data, "")
	item.SetAccessible(keychain.AccessibleWhenUnlockedThisDeviceOnly)
	if aerr := keychain.AddItem(item); aerr != nil {
		return fmt.Errorf("keychain add failed: %w", aerr)
	}
	return nil
}

func Load[T any](key string) (*T, error) {
	data, err := keychain.GetGenericPassword(storeService, key, "", "")
	if err != nil || data == nil {
		if errors.Is(err, keychain.ErrorItemNotFound) || data == nil {
			return nil, ErrItemNotFound
		}
		return nil, fmt.Errorf("keychain query failed: %w", err)
	}

	var obj T
	if err := sonic.Unmarshal(data, &obj); err != nil {
		return nil, fmt.Errorf("failed to unmarshal object %T: %w", obj, err)
	}
	return &obj, nil
}

func Delete(key string) error {
	err := keychain.DeleteGenericPasswordItem(storeService, key)
	if err != nil && !errors.Is(err, keychain.ErrorItemNotFound) {
		return fmt.Errorf("keychain delete failed: %w", err)
	}
	return nil
}
