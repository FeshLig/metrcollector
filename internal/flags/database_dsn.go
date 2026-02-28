package flags

import (
	"errors"
)

type DatabaseDSN string

func (d DatabaseDSN) String() string {
	return string(d)
}

func (d *DatabaseDSN) Set(s string) error {
	if s == "" {
		return errors.New("database dsn is empty")
	}
	*d = DatabaseDSN(s)
	return nil
}
