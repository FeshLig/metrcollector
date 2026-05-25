package audit_test

import (
	"testing"

	"github.com/FeshLig/metrcollector/internal/audit"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var testEventP = audit.Event{
	Timestamp: 1700000000,
	Metrics:   []string{"cpu", "mem"},
	IPAddress: "127.0.0.1",
}

type mockObserver struct {
	events []audit.Event
	err    error
}

func (m *mockObserver) Process(event audit.Event) error {
	m.events = append(m.events, event)
	return m.err
}

func TestPublisher_Notify_SingleObserver(t *testing.T) {
	obs := &mockObserver{}
	p := audit.NewPublisher()
	p.Subscribe(obs)

	p.Notify(testEventP)

	require.Len(t, obs.events, 1)
	assert.Equal(t, testEventP, obs.events[0])
}

func TestPublisher_Notify_MultipleObservers(t *testing.T) {
	obs1, obs2 := &mockObserver{}, &mockObserver{}
	p := audit.NewPublisher()
	p.Subscribe(obs1)
	p.Subscribe(obs2)

	p.Notify(testEventP)

	assert.Len(t, obs1.events, 1)
	assert.Len(t, obs2.events, 1)
}

func TestPublisher_Notify_NoObservers(t *testing.T) {
	p := audit.NewPublisher()
	assert.NotPanics(t, func() { p.Notify(testEventP) })
}
