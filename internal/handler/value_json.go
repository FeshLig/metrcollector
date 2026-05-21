package handler

import (
	"encoding/json"
	"net/http"

	"github.com/FeshLig/metrcollector/internal/dto"
	"github.com/FeshLig/metrcollector/internal/service"
	"github.com/gin-gonic/gin"
)

// ValueJSONHandler handles metric value requests in JSON format.
type ValueJSONHandler struct {
	service service.MetricsService
}

// NewValueJSONHandler creates a new ValueJSONHandler instance.
func NewValueJSONHandler(s service.MetricsService) *ValueJSONHandler {
	return &ValueJSONHandler{
		service: s,
	}
}

// UpdateFromJSON returns metric value in JSON format.
func (h *ValueJSONHandler) UpdateFromJSON(c *gin.Context) {

	var metric dto.Metrics

	c.Writer.Header().Set("Content-Type", "application/json")

	if err := json.NewDecoder(c.Request.Body).Decode(&metric); err != nil {
		c.String(http.StatusBadRequest, "invalid json")
		return
	}

	result, err := h.service.Get(metric)
	if err != nil {
		errCode, msg := HTTPStatusError(err)
		c.String(errCode, msg)
		return
	}

	c.JSON(http.StatusOK, result)

}
