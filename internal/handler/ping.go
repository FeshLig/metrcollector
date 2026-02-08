package handler

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
)

type CheckDB interface {
	Check(ctx context.Context) error
}

type PingHandler struct {
	checker CheckDB
}

func NewPingHandler(checker CheckDB) *PingHandler {
	return &PingHandler{
		checker: checker,
	}
}

func (h *PingHandler) PingPage(c *gin.Context) {

	err := h.checker.Check(c.Request.Context())
	if err != nil {
		c.String(http.StatusInternalServerError, "database ping error")
		return
	}

	c.Status(http.StatusOK)

}
