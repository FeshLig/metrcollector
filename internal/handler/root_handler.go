package handler

import (
	"net/http"

	"github.com/FeshLig/metrcollector/internal/metric"
	"github.com/gin-gonic/gin"
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

func (h *RootHandler) RootPage(c *gin.Context) {

	c.HTML(http.StatusOK, "root.html", gin.H{
		"Gauges":   h.metrics.SnapshotGauges(),
		"Counters": h.metrics.SnapshotCounters(),
	})

	c.Status(http.StatusOK)

}
