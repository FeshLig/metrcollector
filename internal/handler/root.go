package handler

import (
	"net/http"

	"github.com/FeshLig/metrcollector/internal/metric"
	"github.com/FeshLig/metrcollector/internal/service"
	"github.com/gin-gonic/gin"
)

type RootHandler struct {
	service service.MetricsService
}

func NewRootHandler(s service.MetricsService) *RootHandler {
	return &RootHandler{
		service: s,
	}
}

func (h *RootHandler) RootPage(c *gin.Context) {

	gauges := make(map[string]metric.Gauge)
	counters := make(map[string]metric.Counter)

	gaugesMetr := h.service.SnapshotGaugeMetrics()
	countersMetr := h.service.SnapshotCounterMetrics()

	for _, value := range gaugesMetr {
		gauges[value.ID] = metric.Gauge(*value.Value)
	}

	for _, value := range countersMetr {
		counters[value.ID] = metric.Counter(*value.Delta)
	}

	c.HTML(http.StatusOK, "root.html", gin.H{
		"Gauges":   gauges,
		"Counters": counters,
	})

	c.Status(http.StatusOK)

}
