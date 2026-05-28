package middleware_test

import (
	"bytes"
	"compress/gzip"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/FeshLig/metrcollector/internal/middleware"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

func TestGzip_CompressResponse(t *testing.T) {
	gin.SetMode(gin.TestMode)

	r := gin.New()
	r.Use(middleware.Gzip())

	r.GET("/", func(c *gin.Context) {
		c.Header("Content-Type", "application/json")
		c.String(http.StatusOK, `{"status":"ok"}`)
	})

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Accept-Encoding", "gzip")

	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	require.Equal(t, "gzip", w.Header().Get("Content-Encoding"))

	zr, err := gzip.NewReader(bytes.NewReader(w.Body.Bytes()))
	require.NoError(t, err)

	body, err := io.ReadAll(zr)
	require.NoError(t, err)

	require.Equal(t, `{"status":"ok"}`, string(body))
}

func TestGzip_NoCompressResponse(t *testing.T) {
	gin.SetMode(gin.TestMode)

	r := gin.New()
	r.Use(middleware.Gzip())

	r.GET("/", func(c *gin.Context) {
		c.Header("Content-Type", "text/plain")
		c.String(http.StatusOK, "plain text")
	})

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Accept-Encoding", "gzip")

	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	require.Empty(t, w.Header().Get("Content-Encoding"))
	require.Equal(t, "plain text", w.Body.String())
}

func TestGzip_DecompressRequest(t *testing.T) {
	gin.SetMode(gin.TestMode)

	r := gin.New()
	r.Use(middleware.Gzip())

	r.POST("/", func(c *gin.Context) {
		body, err := io.ReadAll(c.Request.Body)
		require.NoError(t, err)

		c.String(http.StatusOK, string(body))
	})

	var buf bytes.Buffer

	zw := gzip.NewWriter(&buf)

	_, err := zw.Write([]byte("compressed body"))
	require.NoError(t, err)

	require.NoError(t, zw.Close())

	req := httptest.NewRequest(http.MethodPost, "/", &buf)
	req.Header.Set("Content-Encoding", "gzip")

	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)
	require.Equal(t, "compressed body", w.Body.String())
}

func TestGzip_InvalidRequestGzip(t *testing.T) {
	gin.SetMode(gin.TestMode)

	r := gin.New()
	r.Use(middleware.Gzip())

	r.POST("/", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	req := httptest.NewRequest(
		http.MethodPost,
		"/",
		bytes.NewBufferString("invalid gzip"),
	)

	req.Header.Set("Content-Encoding", "gzip")

	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusBadRequest, w.Code)
}
