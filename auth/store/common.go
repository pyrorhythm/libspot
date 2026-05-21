package store

import (
	"os"
	"path/filepath"

	"github.com/pyrorhythm/libspot/pkg/keychain"
)

const storeService = "com.pyrorhythm.libspot"

func Zalando[T any](key string) keychain.Keychainer[T] {
	return keychain.NewZalandoKeychainer[T](storeService, key)
}

func File[T any](key string) keychain.Keychainer[T] {
	return keychain.NewFileKeychainer[T](storeDir(), key)
}

func FileCustomDir[T any](key, dir string) keychain.Keychainer[T] {
	return keychain.NewFileKeychainer[T](dir, key)
}

func storeDir() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return ".libspot"
	}
	return filepath.Join(home, ".libspot")
}
