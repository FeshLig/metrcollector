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

func BenchmarkUpdatesHandler_Updates(b *testing.B) {
	gin.SetMode(gin.ReleaseMode)

	storage := repository.NewMemStorage()

	svc := service.NewMetricService(
		storage,
		nil,
		false,
	)

	publisher := audit.NewPublisher()

	handler := NewUpdatesHandler(
		svc,
		publisher,
	)

	router := gin.New()
	router.POST("/updates/", handler.Updates)

	body := []byte(`[
		{
			"id":"Alloc",
			"type":"gauge",
			"value":123.45
		},
		{
			"id":"PollCount",
			"type":"counter",
			"delta":1
		},
		{
			"id":"RandomValue",
			"type":"gauge",
			"value":999.99
		},
		{
			"id":"Requests",
			"type":"counter",
			"delta":42
		}
	]`)

	b.ReportAllocs()

	for b.Loop() {
		b.StopTimer()
		req := httptest.NewRequest(
			"POST",
			"/updates/",
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
