package flags

import (
	"errors"
)

type AuditURL string

func (a AuditURL) String() string {
	return string(a)
}

func (a *AuditURL) Set(s string) error {
	if s == "" {
		return errors.New("audit url is empty")
	}
	*a = AuditURL(s)
	return nil
}
