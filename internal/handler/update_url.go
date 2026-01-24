package handler

import (
	"net/http"
	"strconv"

	"github.com/FeshLig/metrcollector/internal/dto"
	"github.com/FeshLig/metrcollector/internal/service"
	"github.com/gin-gonic/gin"
)

type UpdateURLHandler struct {
	service service.MetricsService
}

func NewUpdateURLHandler(s service.MetricsService) *UpdateURLHandler {
	return &UpdateURLHandler{
		service: s,
	}
}

func (h *UpdateURLHandler) UpdateFromURL(c *gin.Context) {

	metricType := c.Param("type")
	name := c.Param("name")
	valueStr := c.Param("value")

	metric := dto.Metrics{
		ID:    name,
		MType: metricType,
	}

	switch metricType {

	case "gauge":
		value, err := strconv.ParseFloat(valueStr, 64)
		if err != nil {
			c.String(http.StatusBadRequest, "wrong gauge value")
			return
		}
		metric.Value = &value

	case "counter":
		delta, err := strconv.ParseInt(valueStr, 10, 64)
		if err != nil {
			c.String(http.StatusBadRequest, "wrong counter value")
			return
		}
		metric.Delta = &delta

	default:
		c.String(http.StatusBadRequest, "unknown metric type")
		return

	}

	err := h.service.Update(metric)
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	c.Status(http.StatusOK)

}
