package agent_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/FeshLig/metrcollector/internal/agent"
	"github.com/FeshLig/metrcollector/internal/metric"
	"github.com/FeshLig/metrcollector/internal/repository"
)

func TestMetricsCollector_PollCount(t *testing.T) {
	tests := []struct {
		name        string
		iterations  int
		wantCounter metric.Counter
	}{
		{
			name:        "single poll",
			iterations:  1,
			wantCounter: 1,
		},
		{
			name:        "multiple polls",
			iterations:  5,
			wantCounter: 5,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			storage := repository.NewMemStorage()
			collector := agent.NewMetricCollector(storage)

			for i := 0; i < tt.iterations; i++ {
				collector.CollectMetrics()
			}

			counters := storage.SnapshotCounters()

			got, ok := counters["PollCount"]
			require.True(t, ok, "PollCount must exist")
			require.Equal(t, tt.wantCounter, got)
		})
	}
}
