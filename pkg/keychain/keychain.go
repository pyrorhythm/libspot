package keychain

type Error struct {
	msg string
}

func (e *Error) Error() string { return e.msg }
func (e *Error) Unwrap() error { return nil }

var ErrItemNotFound = &Error{"item not found in keychain"}

type Keychainer[T any] interface {
	Load(invalidate bool) (*T, error)
	Save(item *T) error
	Delete() error
}
