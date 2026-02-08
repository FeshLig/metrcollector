package handler

import (
	"net/http"
	"strconv"

	"github.com/FeshLig/metrcollector/internal/dto"
	"github.com/FeshLig/metrcollector/internal/service"
	"github.com/gin-gonic/gin"
)

type ValueURLHandler struct {
	service service.MetricsService
}

func NewValueURLHandler(s service.MetricsService) *ValueURLHandler {
	return &ValueURLHandler{
		service: s,
	}
}

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
	case "gauge":

		c.String(http.StatusOK, strconv.FormatFloat(*result.Value, 'f', -1, 64))

	case "counter":

		c.String(http.StatusOK, strconv.FormatInt(*result.Delta, 10))

	default:
		c.String(http.StatusBadRequest, "unknown metric type")
		return
	}
}
