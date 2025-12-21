package handler

import (
	"github.com/FeshLig/metrcollector/internal/metric"
)

// можно заменить структуру на просто тип мапы
type GaugeHandler struct {
	storage metric.Gauges
}

func NewGaugeHandler(s metric.Gauges) *GaugeHandler {
	return &GaugeHandler{storage: s}
}

func (h *GaugeHandler) Update(name string, value float64) {
	h.storage.SetGauge(name, value)
}
