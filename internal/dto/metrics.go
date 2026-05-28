package dto

const (
	// Gauge identifies gauge metric type.
	Gauge = "gauge"

	// Counter identifies counter metric type.
	Counter = "counter"
)

// Metrics describes metric entity transferred via HTTP API.
type Metrics struct {
	// ID contains metric name.
	ID string `json:"id"`

	// MType contains metric type.
	MType string `json:"type"`

	// Delta contains counter metric value.
	Delta *int64 `json:"delta,omitempty"`

	// Value contains gauge metric value.
	Value *float64 `json:"value,omitempty"`
}
