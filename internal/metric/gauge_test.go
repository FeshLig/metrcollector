package metric_test

import (
	"testing"

	"github.com/FeshLig/metrcollector/internal/metric"
	"github.com/stretchr/testify/assert"
)

func TestGauge_Set(t *testing.T) {
	g := metric.Gauge(1.5)
	assert.InDelta(t, 3.14, float64(g.SetGauge(3.14)), 1e-9)
}

func TestGauge_Immutable(t *testing.T) {
	g := metric.Gauge(1.5)
	g.SetGauge(99.9)
	assert.InDelta(t, 1.5, float64(g), 1e-9)
}
