package handler

import (
	"github.com/FeshLig/metrcollector/internal/metric"
)

// можно заменить структуру на просто тип мапы
type CounterHandler struct {
	storage metric.Counters
}

func NewCounterHandler(s metric.Counters) *CounterHandler {
	return &CounterHandler{storage: s}
}

func (h *CounterHandler) Update(name string, value int64) {
	h.storage.AddCounter(name, value)
}
