package metric_test

import (
	"testing"

	"github.com/FeshLig/metrcollector/internal/metric"
	"github.com/stretchr/testify/assert"
)

func TestCounter_Add(t *testing.T) {
	c := metric.Counter(10)
	assert.Equal(t, metric.Counter(15), c.AddCounter(5))
}

func TestCounter_Add_Negative(t *testing.T) {
	c := metric.Counter(10)
	assert.Equal(t, metric.Counter(7), c.AddCounter(-3))
}

func TestCounter_Set(t *testing.T) {
	c := metric.Counter(10)
	assert.Equal(t, metric.Counter(99), c.SetCounter(99))
}

func TestCounter_Immutable(t *testing.T) {
	c := metric.Counter(10)
	c.AddCounter(5)
	assert.Equal(t, metric.Counter(10), c)
}
