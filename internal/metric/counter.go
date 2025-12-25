package metric

type Counter int64

// type Counters map[string]Counter

func (s Counter) AddCounter(value int64) Counter {
	c := s
	c += Counter(value)
	s = c
	return s
}
