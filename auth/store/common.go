package store

import (
	"github.com/pyrorhythm/libspot/pkg/keychain"
)

const storeService = "com.pyrorhythm.libspot"

func Zalando[T any](key string) keychain.Keychainer[T] {
	return keychain.NewZalandoKeychainer[T](storeService, key)
}

func Keybase[T any](key string) keychain.Keychainer[T] {
	return keychain.NewKeybaseKeychainer[T](storeService, key)
}
