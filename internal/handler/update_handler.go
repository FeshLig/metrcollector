package handler

import (
	"net/http"
	"strconv"
	"strings"
)

type UpdateMetrics interface {
	SetGauge(name string, value float64)
	AddCounter(name string, delta int64)
}

type UpdateHandler struct {
	storage UpdateMetrics
}

func NewUpdateHandler(u UpdateMetrics) *UpdateHandler {
	return &UpdateHandler{
		storage: u,
	}
}

func (h *UpdateHandler) UpdatePage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusBadRequest)
		return
	}

	parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	if len(parts) < 1 || parts[0] != "update" {
		http.Error(w, "invalid path", http.StatusBadRequest)
		return
	} else if len(parts) < 2 {
		http.Error(w, "invalid path", http.StatusBadRequest)
		return
	} else if len(parts) < 3 {
		http.Error(w, "invalid path", http.StatusNotFound)
		return
	} else if len(parts) < 4 {
		http.Error(w, "invalid path", http.StatusBadRequest)
		return
	}

	metricType := parts[1]
	name := parts[2]
	valueStr := parts[3]

	if strings.TrimSpace(name) == "" {
		http.Error(w, "empty name", http.StatusNotFound)
	}

	switch metricType {
	case "gauge":
		value, err := strconv.ParseFloat(valueStr, 64)
		if err != nil {
			http.Error(w, "wrong gauge value", http.StatusBadRequest)
		}
		h.storage.SetGauge(name, value)

	case "counter":
		value, err := strconv.ParseInt(valueStr, 10, 64)
		if err != nil {
			http.Error(w, "wrong counter value", http.StatusBadRequest)
		}
		h.storage.AddCounter(name, value)
	default:
		http.Error(w, "unknown metric type", http.StatusBadRequest)
		return
	}

}
