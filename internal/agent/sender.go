package agent

import (
	"bytes"
	"compress/gzip"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/FeshLig/metrcollector/internal/dto"
	"github.com/FeshLig/metrcollector/internal/repository"
)

type MetricsSender interface {
	SendBatch(metrics []dto.Metrics) error
}

type HTTPSender struct {
	BaseURL string
	Client  *http.Client
}

func NewSender(url string, client *http.Client) *HTTPSender {
	return &HTTPSender{
		BaseURL: url,
		Client:  client,
	}
}

func (h *HTTPSender) Send(metric *dto.Metrics) error {

	body, err := json.Marshal(metric)
	if err != nil {
		return fmt.Errorf("marshal metric: %w", err)
	}

	var buf bytes.Buffer
	gz := gzip.NewWriter(&buf)

	if _, err := gz.Write(body); err != nil {
		return fmt.Errorf("gzip write: %w", err)
	}
	if err := gz.Close(); err != nil {
		return fmt.Errorf("gzip close: %w", err)
	}

	url := fmt.Sprintf("%s/update/", h.BaseURL)
	req, err := http.NewRequest("POST", url, &buf)
	if err != nil {
		return fmt.Errorf("create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Content-Encoding", "gzip")

	resp, err := h.Client.Do(req)
	if err != nil {
		return fmt.Errorf("send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status: %d", resp.StatusCode)
	}

	return nil
}

func (h *HTTPSender) SendBatch(metrics []dto.Metrics) error {

	if len(metrics) == 0 {
		return nil
	}

	body, err := json.Marshal(metrics)
	if err != nil {
		return fmt.Errorf("marshal metric: %w", err)
	}

	var buf bytes.Buffer
	gz := gzip.NewWriter(&buf)

	if _, err := gz.Write(body); err != nil {
		return fmt.Errorf("gzip write: %w", err)
	}
	if err := gz.Close(); err != nil {
		return fmt.Errorf("gzip close: %w", err)
	}

	url := fmt.Sprintf("%s/updates/", h.BaseURL)

	return withRetry(context.TODO(), func() error {

		req, err := http.NewRequest(
			"POST",
			url,
			bytes.NewReader(buf.Bytes()),
		)
		if err != nil {
			return fmt.Errorf("create request: %w", err)
		}

		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Content-Encoding", "gzip")

		resp, err := h.Client.Do(req)
		if err != nil {
			return fmt.Errorf("send request: %w", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode >= 500 {
			return fmt.Errorf("server error: %d", resp.StatusCode)
		}

		if resp.StatusCode != http.StatusOK {
			return fmt.Errorf("bad status: %d", resp.StatusCode)
		}
		return nil
	})

}

func (h *HTTPSender) SendMetrics(storage repository.Storage) {

	sender := h

	for name, value := range storage.SnapshotGauges(context.TODO()) {
		v := float64(value)
		metric := dto.Metrics{
			ID:    name,
			MType: "gauge",
			Value: &v,
		}

		if err := sender.Send(&metric); err != nil {
			log.Printf("ошибка отправки gauge %s: %v", name, err)
		}
	}

	for name, value := range storage.SnapshotCounters(context.TODO()) {
		v := int64(value)
		metric := dto.Metrics{
			ID:    name,
			MType: "counter",
			Delta: &v,
		}
		if err := sender.Send(&metric); err != nil {
			log.Printf("ошибка отправки counter %s: %v", name, err)
		}
	}

}

func (h *HTTPSender) SendBatchMetrics(storage repository.Storage) error {

	gauges := storage.SnapshotGauges(context.TODO())
	counters := storage.SnapshotCounters(context.TODO())

	metrics := make([]dto.Metrics, 0, len(gauges)+len(counters))

	for name, value := range gauges {
		v := float64(value)
		metrics = append(metrics, dto.Metrics{
			ID:    name,
			MType: "gauge",
			Value: &v,
		})
	}

	for name, value := range counters {
		v := int64(value)
		metrics = append(metrics, dto.Metrics{
			ID:    name,
			MType: "counter",
			Delta: &v,
		})
	}

	if err := h.SendBatch(metrics); err != nil {
		return fmt.Errorf("batch send error: %v", err)
	}

	return nil

}

func RunSender(storage *repository.MemStorage, options Options) {
	client := &http.Client{}

	collector := NewMetricCollector(storage)
	sender := NewSender("http://"+options.Address.String(), client)

	pollTicker := time.NewTicker(options.PollInterval.Duration)
	reportTicker := time.NewTicker(options.ReportInterval.Duration)
	defer pollTicker.Stop()
	defer reportTicker.Stop()

	for {

		select {

		case <-pollTicker.C:
			collector.CollectMetrics()

		case <-reportTicker.C:

			if err := sender.SendBatchMetrics(storage); err == nil {
				storage.SetCounter(context.TODO(), "PollCount", 0)
			}

		}
	}

}

func withRetry(ctx context.Context, fn func() error) error {
	const maxRetries = 3
	var lastErr error

	for attempt := 0; attempt <= maxRetries; attempt++ {

		err := fn()
		if err == nil {
			return nil
		}

		if !isRetriable(err) {
			return err
		}

		lastErr = err

		if attempt == maxRetries {
			break
		}

		backoff := time.Duration(1+2*attempt) * time.Second

		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(backoff):
		}
	}

	return lastErr
}

func isRetriable(err error) bool {
	var netErr net.Error
	if errors.As(err, &netErr) {
		return true
	}
	if strings.Contains(err.Error(), "server error") {
		return true
	}

	return false
}
