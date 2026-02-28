package flags

import (
	"errors"
)

type Key string

func (d Key) String() string {
	return string(d)
}

func (d *Key) Set(s string) error {
	if s == "" {
		return errors.New("key is empty")
	}
	*d = Key(s)
	return nil
}
