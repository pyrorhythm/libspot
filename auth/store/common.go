package store

import (
	"github.com/pyrorhythm/libspot/pkg/keychain"
)

const storeService = "com.pyrorhythm.libspot"

var ErrItemNotFound = keychain.ErrItemNotFound

type Keychainer[T any] keychain.Keychainer[T]

func Zalando[T any](key string) Keychainer[T] {
	return keychain.NewZalandoKeychainer[T](storeService, key)
}

func Keybase[T any](key string) Keychainer[T] {
	return keychain.NewKeybaseKeychainer[T](storeService, key)
}
