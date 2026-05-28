package handler

import (
	"bytes"
	"net/http/httptest"
	"testing"

	"github.com/FeshLig/metrcollector/internal/dto"
	"github.com/FeshLig/metrcollector/internal/repository"
	"github.com/FeshLig/metrcollector/internal/service"
	"github.com/gin-gonic/gin"
)

func BenchmarkValueJSONHandler_Gauge(b *testing.B) {
	gin.SetMode(gin.ReleaseMode)

	storage := repository.NewMemStorage()

	svc := service.NewMetricService(
		storage,
		nil,
		false,
	)

	value := 123.45

	err := svc.Update(dto.Metrics{
		ID:    "Alloc",
		MType: dto.Gauge,
		Value: &value,
	})
	if err != nil {
		b.Fatal(err)
	}

	handler := NewValueJSONHandler(svc)

	router := gin.New()
	router.POST("/value/", handler.UpdateFromJSON)

	body := []byte(`{
		"id":"Alloc",
		"type":"gauge"
	}`)

	b.ReportAllocs()

	for b.Loop() {
		req := httptest.NewRequest(
			"POST",
			"/value/",
			bytes.NewReader(body),
		)

		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)
	}
}

func BenchmarkValueJSONHandler_Counter(b *testing.B) {
	gin.SetMode(gin.ReleaseMode)

	storage := repository.NewMemStorage()

	svc := service.NewMetricService(
		storage,
		nil,
		false,
	)

	delta := int64(100)

	err := svc.Update(dto.Metrics{
		ID:    "PollCount",
		MType: dto.Counter,
		Delta: &delta,
	})
	if err != nil {
		b.Fatal(err)
	}

	handler := NewValueJSONHandler(svc)

	router := gin.New()
	router.POST("/value/", handler.UpdateFromJSON)

	body := []byte(`{
		"id":"PollCount",
		"type":"counter"
	}`)

	b.ReportAllocs()

	for b.Loop() {
		req := httptest.NewRequest(
			"POST",
			"/value/",
			bytes.NewReader(body),
		)

		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)
	}
}
