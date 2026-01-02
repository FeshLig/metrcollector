package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ValueHandler struct {
	metrics SnapshotMetrics
}

func NewValueHandler(m SnapshotMetrics) *ValueHandler {
	return &ValueHandler{
		metrics: m,
	}
}

func (h *ValueHandler) ValuePage(c *gin.Context) {

	metricType := c.Param("type")
	name := c.Param("name")

	switch metricType {
	case "gauge":
		gauges := h.metrics.SnapshotGauges()
		value, ok := gauges[name]
		if !ok {
			c.String(http.StatusNotFound, "unknown gauge name")
			return
		}
		c.String(http.StatusOK, strconv.FormatFloat(float64(value), 'f', -1, 64))
	case "counter":
		counters := h.metrics.SnapshotCounters()
		value, ok := counters[name]
		if !ok {
			c.String(http.StatusNotFound, "unknown counter name")
			return
		}
		c.String(http.StatusOK, strconv.FormatInt(int64(value), 10))
	default:
		c.String(http.StatusBadRequest, "unknown metric type")
		return
	}
}
