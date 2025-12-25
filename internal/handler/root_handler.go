package handler

import (
	"fmt"
	"net/http"

	"github.com/FeshLig/metrcollector/internal/repository"
)

type RootHandler struct {
	memStorage *repository.MemStorage
}

func NewRootHandler(m *repository.MemStorage) *RootHandler {
	return &RootHandler{
		memStorage: m,
	}
}

func (h *RootHandler) RootPage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "text/plain")

	fmt.Fprintln(w, "GAUGES:")
	for name, gauge := range h.memStorage.SnapshotGauges() {
		fmt.Fprintf(w, "%s = %f\n", name, gauge)
	}

	fmt.Fprintln(w)
	fmt.Fprintln(w, "COUNTERS:")
	for name, counter := range h.memStorage.SnapshotCounters() {
		fmt.Fprintf(w, "%s = %d\n", name, counter)
	}
}
