package middleware

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
)

type bodyWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w *bodyWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

func Hash(key string) gin.HandlerFunc {
	return func(c *gin.Context) {

		headerHash := c.GetHeader("HashSHA256")
		if key == "" || headerHash == "" {
			c.Next()
			return
		}

		body, err := io.ReadAll(c.Request.Body)
		if err != nil {
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}

		c.Request.Body = io.NopCloser(bytes.NewBuffer(body))

		h := sha256.New()
		h.Write(body)
		h.Write([]byte(key))
		reqHash := hex.EncodeToString(h.Sum(nil))

		if headerHash != reqHash {
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}

		bw := &bodyWriter{
			body:           bytes.NewBuffer(nil),
			ResponseWriter: c.Writer,
		}
		c.Writer = bw

		c.Next()

		h = sha256.New()
		h.Write(bw.body.Bytes())
		h.Write([]byte(key))
		respHash := hex.EncodeToString(h.Sum(nil))

		c.Writer.Header().Set("HashSHA256", respHash)
	}
}
