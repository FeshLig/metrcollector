// Package agent implements metric collection
// and metric delivery to the server.
package agent

import (
	"bytes"
	"compress/gzip"
	"context"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/FeshLig/metrcollector/internal/dto"
	"github.com/FeshLig/metrcollector/internal/repository"
	"github.com/FeshLig/metrcollector/pkg/crypto"
)

// MetricsSender describes metric batch sender.
type MetricsSender interface {
	SendBatch(ctx context.Context, metrics []dto.Metrics, key string, publicKey *rsa.PublicKey) error
}

// Doer describes HTTP client with retry support.
type Doer interface {
	Do(context.Context, func() (*http.Request, error)) (*http.Response, error)
}

// HTTPSender sends metrics to remote HTTP server.
type HTTPSender struct {
	BaseURL string
	Client  Doer
}

// NewSender creates new HTTP metrics sender.
func NewSender(url string, client Doer) *HTTPSender {
	return &HTTPSender{
		BaseURL: url,
		Client:  client,
	}
}

// SendBatch sends metric batch to server.
func (h *HTTPSender) SendBatch(ctx context.Context, metrics []dto.Metrics, key string, publicKey *rsa.PublicKey) error {

	var hash string

	if len(metrics) == 0 {
		return nil
	}

	body, err := json.Marshal(metrics)
	if err != nil {
		return fmt.Errorf("marshal metric: %w", err)
	}

	if key != "" {
		hh := sha256.New()
		hh.Write(body)
		hh.Write([]byte(key))
		hash = hex.EncodeToString(hh.Sum(nil))
	}

	var buf bytes.Buffer
	gz := gzip.NewWriter(&buf)

	if _, err = gz.Write(body); err != nil {
		return fmt.Errorf("gzip write: %w", err)
	}
	if err = gz.Close(); err != nil {
		return fmt.Errorf("gzip close: %w", err)
	}

	payload := buf.Bytes()
	if publicKey != nil {
		payload, err = crypto.Encrypt(publicKey, payload)
		if err != nil {
			return fmt.Errorf("encrypt body: %w", err)
		}
	}

	url := fmt.Sprintf("%s/updates/", h.BaseURL)

	resp, err := h.Client.Do(ctx, func() (*http.Request, error) {
		var req *http.Request
		req, err = http.NewRequest(
			"POST",
			url,
			bytes.NewReader(payload),
		)
		if err != nil {
			return nil, err
		}

		req.Header.Set("Content-Type", "application/json")
		if publicKey == nil {
			req.Header.Set("Content-Encoding", "gzip")
		} else {
			req.Header.Set("X-Encrypted", "rsa-oaep")
		}
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

func buildBatch(storage repository.Storage) []dto.Metrics {
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

	return metrics
}

func worker(
	ctx context.Context,
	jobs <-chan []dto.Metrics,
	sender *HTTPSender,
	key string,
	publicKey *rsa.PublicKey,
	wg *sync.WaitGroup,
) {
	defer wg.Done()
	for {
		select {
		case <-ctx.Done():
			// drain remaining jobs before exit
			for {
				select {
				case metrics, ok := <-jobs:
					if !ok {
						return
					}
					sender.SendBatch(context.Background(), metrics, key, publicKey)
				default:
					return
				}
			}

		case metrics, ok := <-jobs:
			if !ok {
				return
			}
			sender.SendBatch(ctx, metrics, key, publicKey)
		}
	}
}

// RunSender starts metric collection and sending loops.
func RunSender(ctx context.Context, storage *repository.MemStorage, options Options) {

	baseClient := &http.Client{
		Timeout: 5 * time.Second,
	}
	retryClient := NewRetryClient(baseClient)

	sender := NewSender("http://"+options.Address.String(), retryClient)
	collector := NewMetricCollector(storage)
	key := options.Key.String()

	var publicKey *rsa.PublicKey
	if keyPath := options.CryptoKey.String(); keyPath != "" {
		var err error
		publicKey, err = crypto.LoadPublicKey(keyPath)
		if err != nil {
			fmt.Printf("failed to load public key: %v\n", err)
		}
	}

	pollTicker := time.NewTicker(options.PollInterval.Duration)
	reportTicker := time.NewTicker(options.ReportInterval.Duration)
	defer pollTicker.Stop()
	defer reportTicker.Stop()

	jobs := make(chan []dto.Metrics, options.RateLimit)
	collectCh := make(chan struct{})

	var pollCount int64
	var mu sync.Mutex
	var wg sync.WaitGroup

	for i := 0; i < int(options.RateLimit); i++ {
		wg.Add(1)
		go worker(ctx, jobs, sender, key, publicKey, &wg)
	}

	go func() {
		for {
			select {
			case <-ctx.Done():
				close(collectCh)
				return
			case <-pollTicker.C:
				collector.CollectMetrics(ctx)

				mu.Lock()
				pollCount++
				mu.Unlock()

				collectCh <- struct{}{}
			}
		}
	}()

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case <-collectCh:
				collector.CollectGopsutil(ctx)
			}
		}
	}()

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case <-reportTicker.C:
				metrics := buildBatch(storage)
				if len(metrics) > 0 {
					mu.Lock()
					currentPollCount := pollCount
					pollCount = 0
					mu.Unlock()
					metrics = append(metrics, dto.Metrics{
						ID:    "PollCount",
						MType: "counter",
						Delta: &currentPollCount,
					})
					jobs <- metrics
				}
			}
		}
	}()

	<-ctx.Done()

	metrics := buildBatch(storage)
	if len(metrics) > 0 {
		var finalPollCount int64
		mu.Lock()
		finalPollCount = pollCount
		mu.Unlock()
		metrics = append(metrics, dto.Metrics{
			ID:    "PollCount",
			MType: "counter",
			Delta: &finalPollCount,
		})
		jobs <- metrics
	}

	close(jobs)

	wg.Wait()

}
