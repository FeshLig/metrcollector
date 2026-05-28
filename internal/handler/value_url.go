package handler

import (
	"net/http"
	"strconv"

	"github.com/FeshLig/metrcollector/internal/dto"
	"github.com/FeshLig/metrcollector/internal/service"
	"github.com/gin-gonic/gin"
)

// ValueURLHandler handles metric value requests from URL parameters.
type ValueURLHandler struct {
	service service.MetricsService
}

// NewValueURLHandler creates a new ValueURLHandler instance.
func NewValueURLHandler(s service.MetricsService) *ValueURLHandler {
	return &ValueURLHandler{
		service: s,
	}
}

// ValueFromURL returns metric value from URL parameters.
func (h *ValueURLHandler) ValueFromURL(c *gin.Context) {

	name := c.Param("name")
	metricType := c.Param("type")

	metric := dto.Metrics{
		ID:    name,
		MType: metricType,
	}

	result, err := h.service.Get(metric)
	if err != nil {
		c.String(http.StatusNotFound, err.Error())
		return
	}

	switch result.MType {
	case dto.Gauge:

		c.String(http.StatusOK, strconv.FormatFloat(*result.Value, 'f', -1, 64))

	case dto.Counter:

		c.String(http.StatusOK, strconv.FormatInt(*result.Delta, 10))

	default:
		c.String(http.StatusBadRequest, "unknown metric type")
		return
	}
}
