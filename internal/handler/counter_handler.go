package handler

type CounterStorage interface {
	AddCounter(name string, delta int64)
}

type CounterHandler struct {
	storage CounterStorage
}

func NewCounterHandler(s CounterStorage) *CounterHandler {
	return &CounterHandler{storage: s}
}

func (h *CounterHandler) Update(name string, value int64) {
	h.storage.AddCounter(name, value)
}
