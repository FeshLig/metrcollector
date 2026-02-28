package agent

import (
	"bytes"
	"compress/gzip"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/FeshLig/metrcollector/internal/dto"
	"github.com/FeshLig/metrcollector/internal/repository"
)

type MetricsSender interface {
	SendBatch(metrics []dto.Metrics) error
}

type Doer interface {
	Do(context.Context, func() (*http.Request, error)) (*http.Response, error)
}

type HTTPSender struct {
	BaseURL string
	Client  Doer
}

func NewSender(url string, client Doer) *HTTPSender {
	return &HTTPSender{
		BaseURL: url,
		Client:  client,
	}
}

func (h *HTTPSender) SendBatch(ctx context.Context, metrics []dto.Metrics, key string) error {

	var hash string

	if len(metrics) == 0 {
		return nil
	}

	body, err := json.Marshal(metrics)
	if err != nil {
		return fmt.Errorf("marshal metric: %w", err)
	}

	if key != "" {
		h := sha256.New()
		h.Write(body)
		h.Write([]byte(key))
		hash = hex.EncodeToString(h.Sum(nil))
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

	resp, err := h.Client.Do(ctx, func() (*http.Request, error) {

		req, err := http.NewRequest(
			"POST",
			url,
			bytes.NewReader(buf.Bytes()),
		)
		if err != nil {
			return nil, err
		}

		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Content-Encoding", "gzip")
		if key != "" {
			req.Header.Set("HashSHA256", hash)
		}

		return req, nil
	})
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

}

func (h *HTTPSender) SendBatchMetrics(ctx context.Context, storage repository.Storage, key string) error {

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

	if err := h.SendBatch(ctx, metrics, key); err != nil {
		return fmt.Errorf("batch send error: %v", err)
	}

	return nil

}

func RunSender(ctx context.Context, storage *repository.MemStorage, options Options) error {

	baseClient := &http.Client{}
	retryClient := NewRetryClient(baseClient)

	sender := NewSender("http://"+options.Address.String(), retryClient)
	collector := NewMetricCollector(storage)
	key := options.Key.String()

	pollTicker := time.NewTicker(options.PollInterval.Duration)
	reportTicker := time.NewTicker(options.ReportInterval.Duration)

	defer pollTicker.Stop()
	defer reportTicker.Stop()

	for {

		select {

		case <-pollTicker.C:
			collector.CollectMetrics()

		case <-reportTicker.C:

			if err := sender.SendBatchMetrics(ctx, storage, key); err == nil {
				storage.SetCounter(context.TODO(), "PollCount", 0)
			} else {
				return err
			}

		}
	}

}
