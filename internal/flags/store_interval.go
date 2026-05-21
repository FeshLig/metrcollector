package flags

import (
	"fmt"
	"strconv"
	"time"
)

// SecondsStoreInterval represents metric persistence interval.
type SecondsStoreInterval struct {
	time.Duration
}

// String returns interval string representation.
func (i SecondsStoreInterval) String() string {
	return i.Duration.String()
}

// Set parses interval value from string.
func (i *SecondsStoreInterval) Set(s string) error {
	if dur, err := time.ParseDuration(s); err == nil {
		i.Duration = dur
		return nil
	}

	if dur, err := strconv.ParseInt(s, 10, 64); err == nil {
		i.Duration = time.Duration(dur) * time.Second
		return nil
	}

	return fmt.Errorf("invalid duration %q", s)
}
