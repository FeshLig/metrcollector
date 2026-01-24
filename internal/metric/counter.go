package metric

type Counter int64

func (s Counter) AddCounter(value Counter) Counter {
	s += value
	return s
}

func (s Counter) SetCounter(value Counter) Counter {
	return value
}
