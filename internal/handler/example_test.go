package handler_test

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"

	"github.com/FeshLig/metrcollector/internal/audit"
	"github.com/FeshLig/metrcollector/internal/handler"
	"github.com/FeshLig/metrcollector/internal/metric"
	"github.com/FeshLig/metrcollector/internal/repository"
	"github.com/FeshLig/metrcollector/internal/service"
	"github.com/gin-gonic/gin"
)

// ExampleUpdateJSONHandler_UpdateFromJSON demonstrates metric update via JSON request.
func ExampleUpdateJSONHandler_UpdateFromJSON() {

	gin.SetMode(gin.TestMode)

	storage := repository.NewMemStorage()
	svc := service.NewMetricService(storage, nil, false)

	handlers := handler.NewHandlers(svc, audit.NewPublisher())

	router := gin.New()
	router.POST("/update/", handlers.UpdateJSON.UpdateFromJSON)

	body := `{
		"id":"Alloc",
		"type":"gauge",
		"value":123.45
	}`

	req := httptest.NewRequest(
		http.MethodPost,
		"/update/",
		bytes.NewBufferString(body),
	)

	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	fmt.Println(w.Code)

	value, ok := storage.GetGauge(context.TODO(), "Alloc")

	fmt.Println(ok)
	fmt.Println(float64(value))

	// Output:
	// 200
	// true
	// 123.45
}

// ExampleValueJSONHandler_UpdateFromJSON demonstrates metric retrieval via JSON request.
func ExampleValueJSONHandler_UpdateFromJSON() {

	gin.SetMode(gin.TestMode)

	storage := repository.NewMemStorage()

	_ = storage.SetGauge(
		context.TODO(),
		"Alloc",
		metric.Gauge(321.54),
	)

	svc := service.NewMetricService(storage, nil, false)

	handlers := handler.NewHandlers(svc, audit.NewPublisher())

	router := gin.New()
	router.POST("/value/", handlers.ValueJSON.UpdateFromJSON)

	body := `{
		"id":"Alloc",
		"type":"gauge"
	}`

	req := httptest.NewRequest(
		http.MethodPost,
		"/value/",
		bytes.NewBufferString(body),
	)

	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	fmt.Println(w.Code)
	fmt.Println(w.Body.String())

	// Output:
	// 200
	// {"id":"Alloc","type":"gauge","value":321.54}
}

// ExampleUpdateURLHandler_UpdateFromURL demonstrates metric update via URL parameters.
func ExampleUpdateURLHandler_UpdateFromURL() {

	gin.SetMode(gin.TestMode)

	storage := repository.NewMemStorage()
	svc := service.NewMetricService(storage, nil, false)

	handlers := handler.NewHandlers(svc, audit.NewPublisher())

	router := gin.New()
	router.POST(
		"/update/:type/:name/:value/",
		handlers.UpdateURL.UpdateFromURL,
	)

	req := httptest.NewRequest(
		http.MethodPost,
		"/update/counter/PollCount/10/",
		nil,
	)

	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	fmt.Println(w.Code)

	value, ok := storage.GetCounter(context.TODO(), "PollCount")

	fmt.Println(ok)
	fmt.Println(int64(value))

	// Output:
	// 200
	// true
	// 10
}

// ExampleValueURLHandler_ValueFromURL demonstrates metric retrieval via URL parameters.
func ExampleValueURLHandler_ValueFromURL() {

	gin.SetMode(gin.TestMode)

	storage := repository.NewMemStorage()

	_ = storage.SetCounter(
		context.TODO(),
		"PollCount",
		metric.Counter(42),
	)

	svc := service.NewMetricService(storage, nil, false)

	handlers := handler.NewHandlers(svc, audit.NewPublisher())

	router := gin.New()
	router.GET(
		"/value/:type/:name/",
		handlers.ValueURL.ValueFromURL,
	)

	req := httptest.NewRequest(
		http.MethodGet,
		"/value/counter/PollCount/",
		nil,
	)

	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	fmt.Println(w.Code)
	fmt.Println(w.Body.String())

	// Output:
	// 200
	// 42
}
