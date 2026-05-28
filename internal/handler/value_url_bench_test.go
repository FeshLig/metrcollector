package handler

import (
	"net/http/httptest"
	"testing"

	"github.com/FeshLig/metrcollector/internal/dto"
	"github.com/FeshLig/metrcollector/internal/repository"
	"github.com/FeshLig/metrcollector/internal/service"
	"github.com/gin-gonic/gin"
)

func BenchmarkValueURLHandler_Gauge(b *testing.B) {
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

	handler := NewValueURLHandler(svc)

	router := gin.New()
	router.GET("/value/:type/:name", handler.ValueFromURL)

	req := httptest.NewRequest(
		"GET",
		"/value/gauge/Alloc",
		nil,
	)

	b.ReportAllocs()

	for b.Loop() {
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
	}
}

func BenchmarkValueURLHandler_Counter(b *testing.B) {
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

	handler := NewValueURLHandler(svc)

	router := gin.New()
	router.GET("/value/:type/:name", handler.ValueFromURL)

	req := httptest.NewRequest(
		"GET",
		"/value/counter/PollCount",
		nil,
	)

	b.ReportAllocs()

	for b.Loop() {
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
	}
}
