package flags

import (
	"errors"
)

// Restore controls whether metrics should be restored from file.
type Restore bool

// String returns restore flag string representation.
func (r Restore) String() string {

	switch r {
	case true:
		return "true"
	case false:
		return "false"
	}

	return ""
}

// Set parses restore flag value.
func (r *Restore) Set(s string) error {
	switch s {
	case "true":
		*r = true
	case "false":
		*r = false
	default:
		return errors.New("wrong restore value")
	}
	return nil
}
