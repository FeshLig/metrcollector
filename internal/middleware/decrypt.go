package middleware

import (
	"bytes"
	"crypto/rsa"
	"io"
	"net/http"

	"github.com/FeshLig/metrcollector/pkg/crypto"
	"github.com/gin-gonic/gin"
)

// Decrypt decrypts RSA-encrypted request bodies.
// If privateKey is nil, the middleware is a no-op and passes requests through unchanged.
func Decrypt(privateKey *rsa.PrivateKey) gin.HandlerFunc {
	return func(c *gin.Context) {
		if privateKey == nil {
			c.Next()
			return
		}

		encrypted, err := io.ReadAll(c.Request.Body)
		if err != nil {
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}
		defer c.Request.Body.Close()

		decrypted, err := crypto.Decrypt(privateKey, encrypted)
		if err != nil {
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}

		c.Request.Body = io.NopCloser(bytes.NewReader(decrypted))
		c.Request.ContentLength = int64(len(decrypted))

		c.Next()
	}
}
