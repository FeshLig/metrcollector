package flags

import (
	"errors"
)

// AuditFile represents audit log file path.
type AuditFile string

// String returns audit file path.
func (f AuditFile) String() string {
	return string(f)
}

// Set sets audit file path value.
func (f *AuditFile) Set(s string) error {
	if s == "" {
		return errors.New("audit file path is empty")
	}
	*f = AuditFile(s)
	return nil
}
