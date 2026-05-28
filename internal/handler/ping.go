package handler

import (
	"log"
	"net/http"

	"github.com/FeshLig/metrcollector/internal/service"
	"github.com/gin-gonic/gin"
)

// PingHandler handles database ping requests.
type PingHandler struct {
	service service.MetricsService
}

// NewPingHandler creates a new PingHandler instance.
func NewPingHandler(service service.MetricsService) *PingHandler {
	return &PingHandler{
		service: service,
	}
}

// PingPage checks database availability.
func (h *PingHandler) PingPage(c *gin.Context) {

	err := h.service.Check(c.Request.Context())
	if err != nil {
		log.Printf("database ping error: %v\n", err)
		c.String(http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}

	c.Status(http.StatusOK)

}
