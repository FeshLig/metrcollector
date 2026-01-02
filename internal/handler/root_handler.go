package handler

import (
	"net/http"
	"strconv"
	"strings"

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

	gauges := h.metrics.SnapshotGauges()
	counters := h.metrics.SnapshotCounters()

	var str strings.Builder

	for name, value := range gauges {
		str.WriteString(name)
		str.WriteString(": ")
		str.WriteString(strconv.FormatFloat(float64(value), 'f', -1, 64))
		str.WriteByte('\n')
	}

	for name, value := range counters {
		str.WriteString(name)
		str.WriteString(": ")
		str.WriteString(strconv.FormatInt(int64(value), 10))
		str.WriteByte('\n')
	}

	c.String(http.StatusOK, str.String())
}
