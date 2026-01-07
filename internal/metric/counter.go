package metric

type Counter int64

func (s Counter) AddCounter(value int64) Counter {
	c := s
	c += Counter(value)
	s = c
	return s
}
