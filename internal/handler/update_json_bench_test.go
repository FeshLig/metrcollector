package handler

import (
	"bytes"
	"net/http/httptest"
	"testing"

	"github.com/FeshLig/metrcollector/internal/audit"
	"github.com/FeshLig/metrcollector/internal/repository"
	"github.com/FeshLig/metrcollector/internal/service"
	"github.com/gin-gonic/gin"
)

func BenchmarkUpdateJSONHandler_Gauge(b *testing.B) {
	gin.SetMode(gin.ReleaseMode)

	storage := repository.NewMemStorage()

	svc := service.NewMetricService(
		storage,
		nil,
		false,
	)

	publisher := audit.NewPublisher()

	handler := NewUpdateJSONHandler(
		svc,
		publisher,
	)

	router := gin.New()

	router.POST(
		"/update/",
		handler.UpdateFromJSON,
	)

	body := []byte(`{
		"id":"Alloc",
		"type":"gauge",
		"value":123.45
	}`)

	b.ReportAllocs()

	for b.Loop() {
		b.StopTimer()
		req := httptest.NewRequest(
			"POST",
			"/update/",
			bytes.NewReader(body),
		)

		req.Header.Set(
			"Content-Type",
			"application/json",
		)
		b.StartTimer()

		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)
	}
}

func BenchmarkUpdateJSONHandler_Counter(b *testing.B) {
	gin.SetMode(gin.ReleaseMode)

	storage := repository.NewMemStorage()

	svc := service.NewMetricService(
		storage,
		nil,
		false,
	)

	publisher := audit.NewPublisher()

	handler := NewUpdateJSONHandler(
		svc,
		publisher,
	)

	router := gin.New()

	router.POST(
		"/update/",
		handler.UpdateFromJSON,
	)

	body := []byte(`{
		"id":"PollCount",
		"type":"counter",
		"delta":100
	}`)

	b.ReportAllocs()

	for b.Loop() {
		b.StopTimer()
		req := httptest.NewRequest(
			"POST",
			"/update/",
			bytes.NewReader(body),
		)

		req.Header.Set(
			"Content-Type",
			"application/json",
		)
		b.StartTimer()

		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)
	}
}
