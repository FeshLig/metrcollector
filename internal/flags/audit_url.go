package flags

import (
	"errors"
)

// AuditURL represents remote audit endpoint URL.
type AuditURL string

// String returns audit endpoint URL.
func (a AuditURL) String() string {
	return string(a)
}

// Set sets audit endpoint URL value.
func (a *AuditURL) Set(s string) error {
	if s == "" {
		return errors.New("audit url is empty")
	}
	*a = AuditURL(s)
	return nil
}
