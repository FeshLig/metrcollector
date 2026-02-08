package handler

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/FeshLig/metrcollector/internal/dto"
	"github.com/FeshLig/metrcollector/internal/service"
	"github.com/gin-gonic/gin"
)

type ValueJSONHandler struct {
	service service.MetricsService
}

func NewValueJSONHandler(s service.MetricsService) *ValueJSONHandler {
	return &ValueJSONHandler{
		service: s,
	}
}

func (h *ValueJSONHandler) UpdateFromJSON(c *gin.Context) {

	var metric dto.Metrics
	var buf bytes.Buffer

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

	json.NewEncoder(&buf).Encode(result)

	c.String(http.StatusOK, buf.String())

}
