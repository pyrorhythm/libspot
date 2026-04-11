package kchain

import (
	"errors"
	"fmt"

	"github.com/bytedance/sonic"
	"github.com/keybase/go-keychain"
)

const keychainService = "com.pyrorhythm.libspot"

type KeychainError struct {
	msg string
}

func (e *KeychainError) Error() string {
	return e.msg
}
func (e *KeychainError) Unwrap() error { return nil }

var ErrNotFoundInKeychain = &KeychainError{"not found in keychain"}

func Save[T any](key string, obj T) error {
	data, err := sonic.Marshal(obj)
	if err != nil {
		return fmt.Errorf("failed to marshal session: %w", err)
	}

	query := keychain.NewItem()
	query.SetSecClass(keychain.SecClassGenericPassword)
	query.SetService(keychainService)
	query.SetAccount(key)
	query.SetMatchLimit(keychain.MatchLimitOne)
	query.SetReturnData(true)

	_, err = keychain.GetGenericPassword(keychainService, key, "", "")
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
	item := keychain.NewGenericPassword(keychainService, key, "", data, "")
	item.SetAccessible(keychain.AccessibleWhenUnlockedThisDeviceOnly)
	if aerr := keychain.AddItem(item); aerr != nil {
		return fmt.Errorf("keychain add failed: %w", aerr)
	}
	return nil
}

func Load[T any](key string) (*T, error) {
	data, err := keychain.GetGenericPassword(keychainService, key, "", "")
	if err != nil || data == nil {
		if errors.Is(err, keychain.ErrorItemNotFound) || data == nil {
			return nil, ErrNotFoundInKeychain
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
	err := keychain.DeleteGenericPasswordItem(keychainService, key)
	if err != nil && !errors.Is(err, keychain.ErrorItemNotFound) {
		return fmt.Errorf("keychain delete failed: %w", err)
	}
	return nil
}
