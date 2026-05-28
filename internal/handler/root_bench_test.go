package handler

import (
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/FeshLig/metrcollector/internal/dto"
	"github.com/FeshLig/metrcollector/internal/repository"
	"github.com/FeshLig/metrcollector/internal/service"
	"github.com/gin-gonic/gin"
)

func BenchmarkRootHandler_RootPage(b *testing.B) {
	gin.SetMode(gin.ReleaseMode)
	storage := repository.NewMemStorage()
	svc := service.NewMetricService(storage, nil, false)

	// Setup: populate storage with gauge and counter metrics
	for i := range 1000 {
		value := float64(i)
		if err := svc.Update(dto.Metrics{
			ID:    strconv.Itoa(i),
			MType: dto.Gauge,
			Value: &value,
		}); err != nil {
			b.Fatal(err)
		}
	}

	for i := range 1000 {
		delta := int64(i)
		if err := svc.Update(dto.Metrics{
			ID:    "counter_" + strconv.Itoa(i),
			MType: dto.Counter,
			Delta: &delta,
		}); err != nil {
			b.Fatal(err)
		}
	}

	handler := NewRootHandler(svc)
	router := gin.New()
	router.LoadHTMLGlob("../templates/*")
	router.GET("/", handler.RootPage)

	req := httptest.NewRequest("GET", "/", nil)

	// Benchmark: measure RootPage rendering with a populated storage
	b.ReportAllocs()
	for b.Loop() {
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
	}
}
