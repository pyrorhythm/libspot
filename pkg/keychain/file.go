package keychain

import (
	"encoding/base64"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/goccy/go-json"
	"pyrorhythm.dev/fn/errs"
)

var _ Keychainer[any] = (*fileKeychainer[any])(nil)

// fileKeychainer persists credentials to a plain file on disk instead of the
// system keychain. The containing directory is created with 0700 and the file
// with 0600, so secrets stay readable only by the owner.
type fileKeychainer[T any] struct {
	dir    string
	key    string
	cached *T
}

func NewFileKeychainerDirFn[T any](dir string) func(string) Keychainer[T] {
	return func(key string) Keychainer[T] { return &fileKeychainer[T]{dir: dir, key: key} }
}

func NewFileKeychainer[T any](dir, key string) Keychainer[T] {
	return &fileKeychainer[T]{dir: dir, key: key}
}

func (f *fileKeychainer[T]) path() string {
	return filepath.Join(f.dir, f.key+".cred")
}

func (f *fileKeychainer[T]) Load(invalidate bool) (*T, error) {
	if !invalidate && f.cached != nil {
		return f.cached, nil
	}

	data, err := os.ReadFile(f.path())
	if errors.Is(err, os.ErrNotExist) {
		return nil, ErrItemNotFound
	} else if err != nil {
		return nil, fmt.Errorf("failed to read: %w", err)
	}

	b, err := base64.StdEncoding.DecodeString(string(data))
	if err != nil {
		return nil, fmt.Errorf("failed to base64 decode: %w", err)
	}

	f.cached = new(T)
	if err := json.Unmarshal(b, f.cached); err != nil {
		return nil, fmt.Errorf("failed to unmarshal: %w", err)
	}
	return f.cached, nil
}

func (f *fileKeychainer[T]) Save(item *T) error {
	b, err := json.Marshal(item)
	if err != nil {
		return fmt.Errorf("failed to marshal: %w", err)
	}
	data := base64.StdEncoding.EncodeToString(b)

	if err := os.MkdirAll(f.dir, 0o700); err != nil {
		return fmt.Errorf("failed to create dir: %w", err)
	}

	path := f.path()
	if err := os.WriteFile(path, []byte(data), 0o600); err != nil {
		return fmt.Errorf("failed to write: %w", err)
	}
	// Chmod explicitly: WriteFile honors umask and won't tighten perms on a
	// pre-existing file.
	if err := os.Chmod(path, 0o600); err != nil {
		return fmt.Errorf("failed to chmod: %w", err)
	}

	f.cached = item
	return nil
}

func (f *fileKeychainer[T]) Delete() error {
	if err := os.Remove(f.path()); err != nil && !errors.Is(err, os.ErrNotExist) {
		return errs.Wrap(err, "failed to delete")
	}
	f.cached = nil
	return nil
}
