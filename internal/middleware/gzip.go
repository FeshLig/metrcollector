package middleware

import (
	"compress/gzip"
	"io"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// compressWriter compresses HTTP response body using gzip.
type compressWriter struct {
	gin.ResponseWriter
	zw        *gzip.Writer
	modeKnown bool
	useGzip   bool
}

// newCompressWriter creates gzip response writer.
func newCompressWriter(w gin.ResponseWriter) *compressWriter {
	return &compressWriter{
		ResponseWriter: w,
		zw:             gzip.NewWriter(w),
	}
}

// Write writes compressed response body.
func (c *compressWriter) Write(p []byte) (int, error) {
	if !c.modeKnown {
		ct := c.Header().Get("Content-Type")
		c.useGzip =
			strings.Contains(ct, "application/json") ||
				strings.Contains(ct, "text/html")

		if c.useGzip {
			c.Header().Set("Content-Encoding", "gzip")
		}
		c.modeKnown = true
	}

	if c.useGzip {
		return c.zw.Write(p)
	}

	return c.ResponseWriter.Write(p)

}

// Close closes gzip writer.
func (c *compressWriter) Close() error {
	if c.useGzip {
		return c.zw.Close()
	}
	return nil
}

// compressReader decompresses gzip request body.
type compressReader struct {
	io.ReadCloser
	zr *gzip.Reader
}

// newCompressReader creates gzip request reader.
func newCompressReader(r io.ReadCloser) (*compressReader, error) {
	zr, err := gzip.NewReader(r)
	if err != nil {
		return nil, err
	}

	return &compressReader{
		ReadCloser: r,
		zr:         zr,
	}, nil
}

// Read reads decompressed request body.
func (c compressReader) Read(p []byte) (n int, err error) {
	return c.zr.Read(p)
}

// Close closes gzip reader.
func (c *compressReader) Close() error {
	if err := c.zr.Close(); err != nil {
		return err
	}
	return c.ReadCloser.Close()
}

// Gzip compresses HTTP responses and decompresses gzip requests.
func Gzip() gin.HandlerFunc {

	return func(c *gin.Context) {

		acceptEncoding := c.GetHeader("Accept-Encoding")

		if strings.Contains(acceptEncoding, "gzip") {
			cw := newCompressWriter(c.Writer)
			c.Writer = cw
			defer cw.Close()
		}

		contentEncoding := c.GetHeader("Content-Encoding")
		if strings.Contains(contentEncoding, "gzip") {
			cr, err := newCompressReader(c.Request.Body)
			if err != nil {
				c.AbortWithStatus(http.StatusBadRequest)
				return
			}
			c.Request.Body = cr
			defer cr.Close()
		}

		c.Next()

	}

}
