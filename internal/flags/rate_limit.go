package flags

import (
	"fmt"
	"strconv"
)

type RateLimit int64

func (r RateLimit) String() string {
	return strconv.FormatInt(int64(r), 10)
}

func (r *RateLimit) Set(s string) error {

	rate, err := strconv.Atoi(s)
	if err != nil {
		return fmt.Errorf("invalid rate limit %w", err)
	}

	*r = RateLimit(rate)
	return nil
}
