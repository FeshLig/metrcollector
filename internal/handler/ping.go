package handler

import (
	"log"
	"net/http"

	"github.com/FeshLig/metrcollector/internal/service"
	"github.com/gin-gonic/gin"
)

type PingHandler struct {
	service service.MetricsService
}

func NewPingHandler(service service.MetricsService) *PingHandler {
	return &PingHandler{
		service: service,
	}
}

func (h *PingHandler) PingPage(c *gin.Context) {

	err := h.service.Check(c.Request.Context())
	if err != nil {
		log.Printf("database ping error: %v\n", err)
		c.String(http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}

	c.Status(http.StatusOK)

}
