package flags

import (
	"errors"
)

type Restore bool

func (r Restore) String() string {

	switch r {
	case true:
		return "true"
	case false:
		return "false"
	}

	return ""
}

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
