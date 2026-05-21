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

// UpdateJSONHandler handles metric updates from JSON requests.
type UpdateJSONHandler struct {
	service service.MetricsService
	audit   *audit.Publisher
}

// NewUpdateJSONHandler creates a new UpdateJSONHandler instance.
func NewUpdateJSONHandler(s service.MetricsService, audit *audit.Publisher) *UpdateJSONHandler {
	return &UpdateJSONHandler{
		service: s,
		audit:   audit,
	}
}

// UpdateFromJSON updates metric value from JSON request body.
func (h *UpdateJSONHandler) UpdateFromJSON(c *gin.Context) {

	var metric dto.Metrics

	if err := json.NewDecoder(c.Request.Body).Decode(&metric); err != nil {
		c.String(http.StatusBadRequest, "invalid json")
		return
	}

	err := h.service.Update(metric)
	if err != nil {
		errCode, msg := HTTPStatusError(err)
		c.String(errCode, msg)
		return
	}

	event := audit.Event{
		Timestamp: time.Now().Unix(),
		Metrics: []string{
			metric.ID,
		},
		IPAddress: c.ClientIP(),
	}

	h.audit.Notify(event)

	c.Status(http.StatusOK)

}
