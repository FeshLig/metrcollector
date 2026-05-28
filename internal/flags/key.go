package flags

import (
	"errors"
)

// Key represents SHA256 signing key.
type Key string

// String returns key value.
func (d Key) String() string {
	return string(d)
}

// Set sets key value.
func (d *Key) Set(s string) error {
	if s == "" {
		return errors.New("key is empty")
	}
	*d = Key(s)
	return nil
}
