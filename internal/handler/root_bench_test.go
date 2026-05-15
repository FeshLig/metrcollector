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

	svc := service.NewMetricService(
		storage,
		nil,
		false,
	)

	for i := 0; i < 1000; i++ {
		value := float64(i)

		err := svc.Update(dto.Metrics{
			ID:    strconv.Itoa(i),
			MType: dto.Gauge,
			Value: &value,
		})
		if err != nil {
			b.Fatal(err)
		}
	}

	for i := 0; i < 1000; i++ {
		delta := int64(i)

		err := svc.Update(dto.Metrics{
			ID:    "counter_" + strconv.Itoa(i),
			MType: dto.Counter,
			Delta: &delta,
		})
		if err != nil {
			b.Fatal(err)
		}
	}

	handler := NewRootHandler(svc)

	router := gin.New()

	router.LoadHTMLGlob("../templates/*")

	router.GET("/", handler.RootPage)

	req := httptest.NewRequest(
		"GET",
		"/",
		nil,
	)

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
	}
}
