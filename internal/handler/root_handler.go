package handler

import (
	"fmt"
	"net/http"

	"github.com/FeshLig/metrcollector/internal/metric"
)

type SnapshotMetrics interface {
	SnapshotGauges() map[string]metric.Gauge
	SnapshotCounters() map[string]metric.Counter
}

type RootHandler struct {
	metrics SnapshotMetrics
}

func NewRootHandler(m SnapshotMetrics) *RootHandler {
	return &RootHandler{
		metrics: m,
	}
}

func (h *RootHandler) RootPage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "text/plain")

	fmt.Fprintln(w, "GAUGES:")
	for name, gauge := range h.metrics.SnapshotGauges() {
		fmt.Fprintf(w, "%s = %f\n", name, gauge)
	}

	fmt.Fprintln(w)
	fmt.Fprintln(w, "COUNTERS:")
	for name, counter := range h.metrics.SnapshotCounters() {
		fmt.Fprintf(w, "%s = %d\n", name, counter)
	}
}
