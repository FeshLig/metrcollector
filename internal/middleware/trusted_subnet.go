package middleware

import (
	"net"
	"net/http"

	"github.com/FeshLig/metrcollector/internal/flags"
	"github.com/gin-gonic/gin"
)

func TrustedSubnet(subnet flags.TrustedSubnet) gin.HandlerFunc {
	return func(c *gin.Context) {

		if subnet.Subnet == nil {
			c.Next()
			return
		}

		ip := net.ParseIP(c.GetHeader("X-Real-IP"))
		if ip == nil {
			c.AbortWithStatus(http.StatusForbidden)
			return
		}

		if !subnet.Contains(ip) {
			c.AbortWithStatus(http.StatusForbidden)
			return
		}

		c.Next()
	}
}
