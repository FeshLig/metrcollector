package persister_test

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/FeshLig/metrcollector/internal/dto"
	"github.com/FeshLig/metrcollector/internal/metric"
	"github.com/FeshLig/metrcollector/internal/persister"
	"github.com/FeshLig/metrcollector/internal/repository"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newTestPersister(t *testing.T, path string) *persister.FilePersister {
	t.Helper()
	storage := repository.NewMemStorage()
	return persister.NewFilePersister(path, storage, time.Hour)
}

func newTestPersisterWithStorage(t *testing.T, path string, storage repository.Storage) *persister.FilePersister {
	t.Helper()
	return persister.NewFilePersister(path, storage, time.Hour)
}

func tmpPath(t *testing.T) string {
	t.Helper()
	return filepath.Join(t.TempDir(), "metrics.json")
}

func TestFilePersister_Save_CreatesFile(t *testing.T) {
	path := tmpPath(t)
	p := newTestPersister(t, path)

	err := p.Save()
	require.NoError(t, err)
	_, err = os.Stat(path)
	assert.NoError(t, err, "файл должен быть создан")
}

func TestFilePersister_Save_WritesGauge(t *testing.T) {
	path := tmpPath(t)
	storage := repository.NewMemStorage()
	storage.SetGauge(ctx(), "temperature", metric.Gauge(36.6))

	p := newTestPersisterWithStorage(t, path, storage)
	require.NoError(t, p.Save())

	data, err := os.ReadFile(path)
	require.NoError(t, err)

	var metrics []dto.Metrics
	require.NoError(t, json.Unmarshal(data, &metrics))

	require.Len(t, metrics, 1)
	assert.Equal(t, "temperature", metrics[0].ID)
	assert.Equal(t, dto.Gauge, metrics[0].MType)
	require.NotNil(t, metrics[0].Value)
	assert.InDelta(t, 36.6, *metrics[0].Value, 1e-9)
}

func TestFilePersister_Save_WritesCounter(t *testing.T) {
	path := tmpPath(t)
	storage := repository.NewMemStorage()
	storage.SetCounter(ctx(), "requests", metric.Counter(42))

	p := newTestPersisterWithStorage(t, path, storage)
	require.NoError(t, p.Save())

	data, err := os.ReadFile(path)
	require.NoError(t, err)

	var metrics []dto.Metrics
	require.NoError(t, json.Unmarshal(data, &metrics))

	require.Len(t, metrics, 1)
	assert.Equal(t, "requests", metrics[0].ID)
	assert.Equal(t, dto.Counter, metrics[0].MType)
	require.NotNil(t, metrics[0].Delta)
	assert.Equal(t, int64(42), *metrics[0].Delta)
}

func TestFilePersister_Save_EmptyStorage(t *testing.T) {
	path := tmpPath(t)
	p := newTestPersister(t, path)

	require.NoError(t, p.Save())

	data, err := os.ReadFile(path)
	require.NoError(t, err)

	var metrics []dto.Metrics
	require.NoError(t, json.Unmarshal(data, &metrics))
	assert.Empty(t, metrics)
}

func TestFilePersister_Save_CreatesParentDirs(t *testing.T) {
	base := t.TempDir()
	path := filepath.Join(base, "a", "b", "c", "metrics.json")

	p := newTestPersister(t, path)
	require.NoError(t, p.Save())

	_, err := os.Stat(path)
	assert.NoError(t, err)
}

func TestFilePersister_Save_IsAtomic(t *testing.T) {
	path := tmpPath(t)
	p := newTestPersister(t, path)
	require.NoError(t, p.Save())

	_, err := os.Stat(path + ".tmp")
	assert.True(t, os.IsNotExist(err), ".tmp файл не должен оставаться после Save")
}

func TestFilePersister_Load_RestoresGauge(t *testing.T) {
	path := tmpPath(t)
	storage := repository.NewMemStorage()
	storage.SetGauge(ctx(), "cpu", metric.Gauge(0.75))

	p1 := newTestPersisterWithStorage(t, path, storage)
	require.NoError(t, p1.Save())

	storage2 := repository.NewMemStorage()
	p2 := newTestPersisterWithStorage(t, path, storage2)
	require.NoError(t, p2.Load())

	gauges, _ := storage2.SnapshotMetrics(ctx())
	val, ok := gauges["cpu"]
	require.True(t, ok)
	assert.InDelta(t, 0.75, float64(val), 1e-9)
}

func TestFilePersister_Load_RestoresCounter(t *testing.T) {
	path := tmpPath(t)
	storage := repository.NewMemStorage()
	storage.SetCounter(ctx(), "hits", metric.Counter(100))

	p1 := newTestPersisterWithStorage(t, path, storage)
	require.NoError(t, p1.Save())

	storage2 := repository.NewMemStorage()
	p2 := newTestPersisterWithStorage(t, path, storage2)
	require.NoError(t, p2.Load())

	_, counters := storage2.SnapshotMetrics(ctx())
	val, ok := counters["hits"]
	require.True(t, ok)
	assert.Equal(t, metric.Counter(100), val)
}

func TestFilePersister_Load_FileNotExist(t *testing.T) {
	path := filepath.Join(t.TempDir(), "nonexistent.json")
	p := newTestPersister(t, path)

	err := p.Load()
	assert.NoError(t, err)
}

func TestFilePersister_Load_InvalidJSON(t *testing.T) {
	path := tmpPath(t)
	require.NoError(t, os.WriteFile(path, []byte("not valid json {{"), 0644))

	p := newTestPersister(t, path)
	err := p.Load()
	assert.Error(t, err)
}

func TestFilePersister_Load_GaugeNilValue(t *testing.T) {
	path := tmpPath(t)
	metrics := []dto.Metrics{
		{ID: "broken", MType: dto.Gauge, Value: nil},
	}
	data, _ := json.Marshal(metrics)
	require.NoError(t, os.WriteFile(path, data, 0644))

	p := newTestPersister(t, path)
	err := p.Load()
	assert.Error(t, err)
}

func TestFilePersister_Load_CounterNilDelta(t *testing.T) {
	path := tmpPath(t)
	metrics := []dto.Metrics{
		{ID: "broken", MType: dto.Counter, Delta: nil},
	}
	data, _ := json.Marshal(metrics)
	require.NoError(t, os.WriteFile(path, data, 0644))

	p := newTestPersister(t, path)
	err := p.Load()
	assert.Error(t, err)
}

func TestFilePersister_Load_UnknownMetricType(t *testing.T) {
	path := tmpPath(t)
	metrics := []dto.Metrics{
		{ID: "x", MType: "histogram"},
	}
	data, _ := json.Marshal(metrics)
	require.NoError(t, os.WriteFile(path, data, 0644))

	p := newTestPersister(t, path)
	err := p.Load()
	assert.Error(t, err)
}

func TestFilePersister_RoundTrip(t *testing.T) {
	path := tmpPath(t)
	storage := repository.NewMemStorage()
	storage.SetGauge(ctx(), "mem", metric.Gauge(1024.5))
	storage.SetCounter(ctx(), "reqs", metric.Counter(7))

	p1 := newTestPersisterWithStorage(t, path, storage)
	require.NoError(t, p1.Save())

	storage2 := repository.NewMemStorage()
	p2 := newTestPersisterWithStorage(t, path, storage2)
	require.NoError(t, p2.Load())

	gauges, counters := storage2.SnapshotMetrics(ctx())

	g, ok := gauges["mem"]
	require.True(t, ok)
	assert.InDelta(t, 1024.5, float64(g), 1e-9)

	c, ok := counters["reqs"]
	require.True(t, ok)
	assert.Equal(t, metric.Counter(7), c)
}

func TestFilePersister_Start_WritesOnTick(t *testing.T) {
	path := tmpPath(t)
	storage := repository.NewMemStorage()
	storage.SetCounter(ctx(), "tick_test", metric.Counter(1))

	p := persister.NewFilePersister(path, storage, 50*time.Millisecond)
	p.Start()
	defer p.Stop()

	time.Sleep(150 * time.Millisecond)

	_, err := os.Stat(path)
	assert.NoError(t, err, "файл должен появиться после тика")
}

func TestFilePersister_Start_ZeroInterval_DoesNotPanic(t *testing.T) {
	path := tmpPath(t)
	p := persister.NewFilePersister(path, repository.NewMemStorage(), 0)
	assert.NotPanics(t, func() {
		p.Start()
	})
}

func TestFilePersister_SaveNow(t *testing.T) {
	path := tmpPath(t)
	storage := repository.NewMemStorage()
	storage.SetGauge(ctx(), "instant", metric.Gauge(9.9))

	p := newTestPersisterWithStorage(t, path, storage)
	p.SaveNow()

	_, err := os.Stat(path)
	assert.NoError(t, err)
}

func ctx() context.Context {
	return context.Background()
}
