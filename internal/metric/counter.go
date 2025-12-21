package metric

type Counter int64

type Counters map[string]Counter

func (s Counters) AddCounter(name string, value int64) {
	c := s[name]
	c += Counter(value)
	s[name] = c
}
