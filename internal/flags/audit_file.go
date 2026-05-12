package flags

import (
	"errors"
)

type AuditFile string

func (f AuditFile) String() string {
	return string(f)
}

func (f *AuditFile) Set(s string) error {
	if s == "" {
		return errors.New("audit file path is empty")
	}
	*f = AuditFile(s)
	return nil
}
