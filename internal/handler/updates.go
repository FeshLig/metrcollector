package handler

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/FeshLig/metrcollector/internal/audit"
	"github.com/FeshLig/metrcollector/internal/dto"
	"github.com/FeshLig/metrcollector/internal/service"
	"github.com/gin-gonic/gin"
)

// UpdatesHandler handles batch metric updates.
type UpdatesHandler struct {
	service service.MetricsService
	audit   *audit.Publisher
}

// NewUpdatesHandler creates a new UpdatesHandler instance.
func NewUpdatesHandler(s service.MetricsService, audit *audit.Publisher) *UpdatesHandler {
	return &UpdatesHandler{
		service: s,
		audit:   audit,
	}
}

// Updates updates multiple metrics from JSON request body.
func (h *UpdatesHandler) Updates(c *gin.Context) {

	var metrics []dto.Metrics

	if err := json.NewDecoder(c.Request.Body).Decode(&metrics); err != nil {
		c.String(http.StatusBadRequest, http.StatusText(http.StatusBadRequest))
		return
	}

	if len(metrics) == 0 {
		c.Status(http.StatusOK)
		return
	}

	if err := h.service.Updates(metrics); err != nil {
		errCode, msg := HTTPStatusError(err)
		c.String(errCode, msg)
		return
	}

	metricNames := make([]string, 0, len(metrics))
	for _, metric := range metrics {
		metricNames = append(metricNames, metric.ID)
	}
	event := audit.Event{
		Timestamp: time.Now().Unix(),
		Metrics:   metricNames,
		IPAddress: c.ClientIP(),
	}

	h.audit.Notify(event)

	c.Status(http.StatusOK)

}
