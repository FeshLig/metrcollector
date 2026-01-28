package middleware

import (
	"compress/gzip"
	"io"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type compressWriter struct {
	gin.ResponseWriter
	zw        *gzip.Writer
	modeKnown bool
	useGzip   bool
}

func newCompressWriter(w gin.ResponseWriter) *compressWriter {
	return &compressWriter{
		ResponseWriter: w,
		zw:             gzip.NewWriter(w),
	}
}

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

func (c *compressWriter) Close() error {
	if c.useGzip {
		return c.zw.Close()
	}
	return nil
}

type compressReader struct {
	io.ReadCloser
	zr *gzip.Reader
}

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

func (c compressReader) Read(p []byte) (n int, err error) {
	return c.zr.Read(p)
}

func (c *compressReader) Close() error {
	if err := c.zr.Close(); err != nil {
		return err
	}
	return c.ReadCloser.Close()
}

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
