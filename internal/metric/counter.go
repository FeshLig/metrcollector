package metric

type Counter int64

func (s Counter) AddCounter(value int64) Counter {
	s += Counter(value)
	return s
}

func (s Counter) SetCounter(value int64) Counter {
	return Counter(value)
}
