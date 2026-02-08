package persister

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"time"

	"github.com/FeshLig/metrcollector/internal/dto"
	"github.com/FeshLig/metrcollector/internal/metric"
)

type MetricsStorage interface {
	SnapshotMetrics() (map[string]metric.Gauge, map[string]metric.Counter)
	SetMetrics(map[string]metric.Gauge, map[string]metric.Counter)
}

type FilePersister struct {
	path     string
	storage  MetricsStorage
	interval time.Duration
	stop     chan struct{}
}

func NewFilePersister(
	path string,
	storage MetricsStorage,
	interval time.Duration,
) *FilePersister {
	return &FilePersister{
		path:     path,
		storage:  storage,
		interval: interval,
		stop:     make(chan struct{}),
	}
}

func (p *FilePersister) Save() error {

	var metrics []dto.Metrics
	gauges, counters := p.storage.SnapshotMetrics()

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
	p.storage.SetMetrics(gauges, counters)

	return nil

}

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

func (p *FilePersister) Stop() {
	close(p.stop)
}

func (p *FilePersister) SaveNow() {
	_ = p.Save()
}
