package handler

import (
	"net/http/httptest"
	"testing"

	"github.com/FeshLig/metrcollector/internal/audit"
	"github.com/FeshLig/metrcollector/internal/repository"
	"github.com/FeshLig/metrcollector/internal/service"
	"github.com/gin-gonic/gin"
)

func BenchmarkUpdateURLHandler_Gauge(b *testing.B) {
	gin.SetMode(gin.ReleaseMode)

	storage := repository.NewMemStorage()

	svc := service.NewMetricService(
		storage,
		nil,
		false,
	)

	publisher := audit.NewPublisher()

	handler := NewUpdateURLHandler(
		svc,
		publisher,
	)

	router := gin.New()

	router.POST(
		"/update/:type/:name/:value",
		handler.UpdateFromURL,
	)

	req := httptest.NewRequest(
		"POST",
		"/update/gauge/Alloc/123.45",
		nil,
	)

	b.ReportAllocs()

	for b.Loop() {
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
	}
}

func BenchmarkUpdateURLHandler_Counter(b *testing.B) {
	gin.SetMode(gin.ReleaseMode)

	storage := repository.NewMemStorage()

	svc := service.NewMetricService(
		storage,
		nil,
		false,
	)

	publisher := audit.NewPublisher()

	handler := NewUpdateURLHandler(
		svc,
		publisher,
	)

	router := gin.New()

	router.POST(
		"/update/:type/:name/:value",
		handler.UpdateFromURL,
	)

	req := httptest.NewRequest(
		"POST",
		"/update/counter/PollCount/100",
		nil,
	)

	b.ReportAllocs()

	for b.Loop() {
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
	}
}
