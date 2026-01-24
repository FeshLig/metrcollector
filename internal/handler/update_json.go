package handler

import (
	"encoding/json"
	"net/http"

	"github.com/FeshLig/metrcollector/internal/dto"
	"github.com/FeshLig/metrcollector/internal/service"
	"github.com/gin-gonic/gin"
)

type UpdateJSONHandler struct {
	service service.MetricsService
}

func NewUpdateJSONHandler(s service.MetricsService) *UpdateJSONHandler {
	return &UpdateJSONHandler{
		service: s,
	}
}

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

	c.Status(http.StatusOK)

}
