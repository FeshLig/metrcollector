package handler

type GaugeStorage interface {
	SetGauge(name string, value float64)
}

// можно заменить структуру на просто тип мапы
type GaugeHandler struct {
	storage GaugeStorage
}

func NewGaugeHandler(s GaugeStorage) *GaugeHandler {
	return &GaugeHandler{storage: s}
}

func (h *GaugeHandler) Update(name string, value float64) {
	h.storage.SetGauge(name, value)
}
