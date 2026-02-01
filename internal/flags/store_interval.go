package flags

import (
	"fmt"
	"strconv"
	"time"
)

type SecondsStoreInterval struct {
	time.Duration
}

func (i SecondsStoreInterval) String() string {
	return i.Duration.String()
}

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
