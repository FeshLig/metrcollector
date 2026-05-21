package flags

import (
	"fmt"
	"strconv"
	"time"
)

// SecondsDuration represents duration value parsed from flags.
type SecondsDuration struct {
	time.Duration
}

// String returns duration string representation.
func (d SecondsDuration) String() string {
	return d.Duration.String()
}

// Set parses duration value from string.
func (d *SecondsDuration) Set(s string) error {
	if dur, err := time.ParseDuration(s); err == nil {
		d.Duration = dur
		return nil
	}

	if dur, err := strconv.ParseInt(s, 10, 64); err == nil {
		d.Duration = time.Duration(dur) * time.Second
		return nil
	}

	return fmt.Errorf("invalid duration %q", s)
}
