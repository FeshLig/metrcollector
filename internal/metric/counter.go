// Package metric provides metric types
// and metric-related operations.
package metric

// Counter represents counter metric value.
type Counter int64

// AddCounter increments counter metric value.
func (s Counter) AddCounter(value Counter) Counter {
	s += value
	return s
}

// SetCounter sets counter metric value.
func (s Counter) SetCounter(value Counter) Counter {
	return value
}
