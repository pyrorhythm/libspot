package store

import (
	"os"
	"path/filepath"

	"pyrorhythm.dev/libspot/pkg/keychain"
)

const storeService = "com.pyrorhythm.libspot"

func Zalando[T any](key string) keychain.Keychainer[T] {
	return keychain.NewZalandoKeychainer[T](storeService, key)
}

func File[T any](key string) keychain.Keychainer[T] {
	return keychain.NewFileKeychainer[T](storeDir(), key)
}

func storeDir() string {
	envdir := os.Getenv("LIBSPOT_ROOT")
	if envdir != "" {
		return envdir
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return ".libspot"
	}
	return filepath.Join(home, ".libspot")
}
