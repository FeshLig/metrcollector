package metric

type MemStorage struct {
	gauges   Gauges
	counters Counters
}

// Сделать NewGauges и NewCounters
func NewMemStorage() *MemStorage {
	return &MemStorage{
		gauges:   make(map[string]Gauge),
		counters: make(map[string]Counter),
	}
}

func (s *MemStorage) GetGauges() Gauges {
	return s.gauges
}

func (s *MemStorage) GetCounters() Counters {
	return s.counters
}
