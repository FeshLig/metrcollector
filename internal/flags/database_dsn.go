package flags

import (
	"errors"
)

// DatabaseDSN represents PostgreSQL connection string.
type DatabaseDSN string

// String returns database connection string.
func (d DatabaseDSN) String() string {
	return string(d)
}

// Set sets database connection string value.
func (d *DatabaseDSN) Set(s string) error {
	if s == "" {
		return errors.New("database dsn is empty")
	}
	*d = DatabaseDSN(s)
	return nil
}
