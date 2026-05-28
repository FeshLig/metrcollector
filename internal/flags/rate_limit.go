package flags

import (
	"fmt"
	"strconv"
)

// RateLimit represents maximum number of concurrent requests.
type RateLimit int64

// String returns rate limit string representation.
func (r RateLimit) String() string {
	return strconv.FormatInt(int64(r), 10)
}

// Set parses rate limit value from string.
func (r *RateLimit) Set(s string) error {

	rate, err := strconv.Atoi(s)
	if err != nil {
		return fmt.Errorf("invalid rate limit %w", err)
	}

	*r = RateLimit(rate)
	return nil
}
