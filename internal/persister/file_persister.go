// Package persister provides file persistence
// for application metrics.
package persister

import (
	"context"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"time"

	"github.com/FeshLig/metrcollector/internal/dto"
	"github.com/FeshLig/metrcollector/internal/metric"
	"github.com/FeshLig/metrcollector/internal/repository"
)

// FilePersister periodically saves metrics into file storage.
type FilePersister struct {
	path     string
	storage  repository.Storage
	interval time.Duration
	stop     chan struct{}
}

// NewFilePersister creates new file persister instance.
func NewFilePersister(
	path string,
	storage repository.Storage,
	interval time.Duration,
) *FilePersister {
	return &FilePersister{
		path:     path,
		storage:  storage,
		interval: interval,
		stop:     make(chan struct{}),
	}
}

// Save writes all metrics into file storage.
func (p *FilePersister) Save() error {

	var metrics []dto.Metrics
	gauges, counters := p.storage.SnapshotMetrics(context.TODO())

	for name, value := range gauges {
		v := float64(value)
		metrics = append(metrics, dto.Metrics{
			ID:    name,
			MType: dto.Gauge,
			Value: &v,
		})
	}
	for name, value := range counters {
		v := int64(value)
		metrics = append(metrics, dto.Metrics{
			ID:    name,
			MType: dto.Counter,
			Delta: &v,
		})
	}

	data, err := json.MarshalIndent(metrics, "", "  ")
	if err != nil {
		return err
	}

	dir := filepath.Dir(p.path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	tmpFile := p.path + ".tmp"
	if err := os.WriteFile(tmpFile, data, 0644); err != nil {
		return err
	}

	return os.Rename(tmpFile, p.path)

}

// Load restores metrics from file storage.
func (p *FilePersister) Load() error {

	var metrics []dto.Metrics
	gauges := make(map[string]metric.Gauge)
	counters := make(map[string]metric.Counter)

	data, err := os.ReadFile(p.path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil
		}
		return err
	}

	if err := json.Unmarshal(data, &metrics); err != nil {
		return err
	}

	for _, m := range metrics {
		switch m.MType {
		case dto.Gauge:
			if m.Value == nil {
				return errors.New("gauge value is nil")
			}
			gauges[m.ID] = metric.Gauge(*m.Value)
		case dto.Counter:
			if m.Delta == nil {
				return errors.New("counter delta is nil")
			}
			counters[m.ID] = metric.Counter(*m.Delta)
		default:
			return errors.New("invalid metric type")
		}
	}
	p.storage.SetMetrics(context.TODO(), gauges, counters)

	return nil

}

// Start launches periodic metric persistence.
func (p *FilePersister) Start() {
	if p.interval <= 0 {
		return
	}

	go func() {
		ticker := time.NewTicker(p.interval)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				_ = p.Save()
			case <-p.stop:
				return
			}
		}
	}()
}

// Stop stops periodic metric persistence.
func (p *FilePersister) Stop() {
	close(p.stop)
}

// SaveNow immediately saves metrics into file storage.
func (p *FilePersister) SaveNow() {
	_ = p.Save()
}
