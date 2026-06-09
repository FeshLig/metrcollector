package router_test

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/FeshLig/metrcollector/internal/audit"
	"github.com/FeshLig/metrcollector/internal/config"
	"github.com/FeshLig/metrcollector/internal/handler"
	"github.com/FeshLig/metrcollector/internal/persister"
	"github.com/FeshLig/metrcollector/internal/repository"
	"github.com/FeshLig/metrcollector/internal/router"
	"github.com/FeshLig/metrcollector/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func setupRouter(t *testing.T) *gin.Engine {
	t.Helper()

	storage := repository.NewMemStorage()

	persister := persister.NewFilePersister(
		"",
		storage,
		time.Hour,
	)

	svc := service.NewMetricService(
		storage,
		persister,
		false,
	)

	auditPublisher := audit.NewPublisher()

	handlers := handler.NewHandlers(
		svc,
		auditPublisher,
	)

	cfg := config.Options{}

	logger := zap.NewNop()

	return router.NewRouter(
		handlers,
		cfg,
		logger,
		nil,
	)
}

func TestRouter_NotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)

	r := setupRouter(t)

	req := httptest.NewRequest(
		http.MethodGet,
		"/missing",
		nil,
	)

	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusNotFound, w.Code)
}

func TestRouter_UpdateRoute(t *testing.T) {
	gin.SetMode(gin.TestMode)

	r := setupRouter(t)

	req := httptest.NewRequest(
		http.MethodPost,
		"/update/counter/test/10/",
		nil,
	)

	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	require.NotEqual(t, http.StatusNotFound, w.Code)
}

func TestRouter_ValueRoute(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := setupRouter(t)

	update := httptest.NewRequest(http.MethodPost, "/update/counter/test/10/", nil)
	httptest.NewRecorder()
	uw := httptest.NewRecorder()
	r.ServeHTTP(uw, update)
	require.Equal(t, http.StatusOK, uw.Code)

	req := httptest.NewRequest(http.MethodGet, "/value/counter/test/", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	require.NotEqual(t, http.StatusNotFound, w.Code)
}

func TestRouter_UsesGzipMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)

	r := setupRouter(t)

	req := httptest.NewRequest(
		http.MethodGet,
		"/",
		nil,
	)

	req.Header.Set("Accept-Encoding", "gzip")

	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	require.Equal(
		t,
		"gzip",
		w.Header().Get("Content-Encoding"),
	)
}

func TestRouter_MethodNotAllowed(t *testing.T) {
	gin.SetMode(gin.TestMode)

	r := setupRouter(t)

	req := httptest.NewRequest(
		http.MethodDelete,
		"/update/",
		nil,
	)

	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	require.True(
		t,
		w.Code == http.StatusNotFound ||
			w.Code == http.StatusMethodNotAllowed,
	)
}
