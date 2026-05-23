package middleware_test

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/FeshLig/metrcollector/internal/middleware"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

func makeHash(t *testing.T, body []byte, key string) string {
	t.Helper()
	h := sha256.New()
	h.Write(body)
	h.Write([]byte(key))

	return hex.EncodeToString(h.Sum(nil))
}

func TestHash_ValidRequest(t *testing.T) {
	gin.SetMode(gin.TestMode)

	key := "secret"

	r := gin.New()
	r.Use(middleware.Hash(key))

	r.POST("/", func(c *gin.Context) {
		body, err := io.ReadAll(c.Request.Body)
		require.NoError(t, err)

		c.String(http.StatusOK, string(body))
	})

	body := []byte("hello")

	hash := makeHash(t, body, key)

	req := httptest.NewRequest(
		http.MethodPost,
		"/",
		bytes.NewBuffer(body),
	)

	req.Header.Set("HashSHA256", hash)

	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)
	require.Equal(t, "hello", w.Body.String())
}

func TestHash_InvalidRequest(t *testing.T) {
	gin.SetMode(gin.TestMode)

	r := gin.New()
	r.Use(middleware.Hash("secret"))

	r.POST("/", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	req := httptest.NewRequest(
		http.MethodPost,
		"/",
		bytes.NewBufferString("hello"),
	)

	req.Header.Set("HashSHA256", "invalid")

	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusBadRequest, w.Code)
}

func TestHash_ResponseHash(t *testing.T) {
	gin.SetMode(gin.TestMode)

	key := "secret"

	r := gin.New()
	r.Use(middleware.Hash(key))

	r.POST("/", func(c *gin.Context) {
		c.String(http.StatusOK, "response")
	})

	body := []byte("request")

	reqHash := makeHash(t, body, key)

	req := httptest.NewRequest(
		http.MethodPost,
		"/",
		bytes.NewBuffer(body),
	)

	req.Header.Set("HashSHA256", reqHash)

	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	expectedRespHash := makeHash(t, []byte("response"), key)

	require.Equal(
		t,
		expectedRespHash,
		w.Header().Get("HashSHA256"),
	)
}

func TestHash_EmptyKey(t *testing.T) {
	gin.SetMode(gin.TestMode)

	r := gin.New()
	r.Use(middleware.Hash(""))

	r.POST("/", func(c *gin.Context) {
		c.String(http.StatusOK, "ok")
	})

	req := httptest.NewRequest(
		http.MethodPost,
		"/",
		bytes.NewBufferString("body"),
	)

	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)
	require.Empty(t, w.Header().Get("HashSHA256"))
}
