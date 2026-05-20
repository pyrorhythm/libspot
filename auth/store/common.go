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

// File stores credentials in ~/.libspot as an alternative to the system
// keychain, useful on headless hosts where no keychain is available.
func File[T any](key string) keychain.Keychainer[T] {
	return keychain.NewFileKeychainer[T](storeDir(), key)
}

func storeDir() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return ".libspot"
	}
	return filepath.Join(home, ".libspot")
}
