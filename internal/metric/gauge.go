package metric

// Gauge represents gauge metric value.
type Gauge float64

// SetGauge sets gauge metric value.
func (s Gauge) SetGauge(value Gauge) Gauge {
	s = value
	return s
}
