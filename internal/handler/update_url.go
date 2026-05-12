package handler

import (
	"net/http"
	"strconv"
	"time"

	"github.com/FeshLig/metrcollector/internal/audit"
	"github.com/FeshLig/metrcollector/internal/dto"
	"github.com/FeshLig/metrcollector/internal/service"
	"github.com/gin-gonic/gin"
)

type UpdateURLHandler struct {
	service service.MetricsService
	audit   *audit.Publisher
}

func NewUpdateURLHandler(s service.MetricsService, audit *audit.Publisher) *UpdateURLHandler {
	return &UpdateURLHandler{
		service: s,
		audit:   audit,
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

	case dto.Gauge:
		value, err := strconv.ParseFloat(valueStr, 64)
		if err != nil {
			c.String(http.StatusBadRequest, "wrong gauge value")
			return
		}
		metric.Value = &value

	case dto.Counter:
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
