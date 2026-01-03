package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type UpdateMetrics interface {
	SetGauge(name string, value float64)
	AddCounter(name string, delta int64)
}

type UpdateHandler struct {
	storage UpdateMetrics
}

func NewUpdateHandler(u UpdateMetrics) *UpdateHandler {
	return &UpdateHandler{
		storage: u,
	}
}

func (h *UpdateHandler) UpdatePage(c *gin.Context) {

	metricType := c.Param("type")
	name := c.Param("name")
	valueStr := c.Param("value")

	// TODO: Возможно перенести в сервис
	switch metricType {
	case "gauge":
		value, err := strconv.ParseFloat(valueStr, 64)
		if err != nil {
			c.String(http.StatusBadRequest, "wrong gauge value")
			return
		}
		h.storage.SetGauge(name, value)

	case "counter":
		value, err := strconv.ParseInt(valueStr, 10, 64)
		if err != nil {
			c.String(http.StatusBadRequest, "wrong counter value")
			return
		}
		h.storage.AddCounter(name, value)
	default:
		c.String(http.StatusBadRequest, "unknown metric type")
		return
	}

	c.Status(http.StatusOK)

}
