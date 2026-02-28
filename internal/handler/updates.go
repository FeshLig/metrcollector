package handler

import (
	"encoding/json"
	"net/http"

	"github.com/FeshLig/metrcollector/internal/dto"
	"github.com/FeshLig/metrcollector/internal/service"
	"github.com/gin-gonic/gin"
)

type UpdatesHandler struct {
	service service.MetricsService
}

func NewUpdatesHandler(s service.MetricsService) *UpdatesHandler {
	return &UpdatesHandler{
		service: s,
	}
}

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

	c.Status(http.StatusOK)

}
