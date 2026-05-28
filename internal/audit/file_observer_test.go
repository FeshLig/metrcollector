package audit_test

import (
	"bufio"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/FeshLig/metrcollector/internal/audit"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var testEventFO = audit.Event{
	Timestamp: 1700000000,
	Metrics:   []string{"cpu", "mem"},
	IPAddress: "127.0.0.1",
}

func TestFileObserver_WritesEvent(t *testing.T) {
	path := filepath.Join(t.TempDir(), "audit.log")
	obs, err := audit.NewFileObserver(path)
	require.NoError(t, err)
	defer obs.Close()

	require.NoError(t, obs.Process(testEventFO))

	f, err := os.Open(path)
	require.NoError(t, err)
	defer f.Close()

	var got audit.Event
	require.NoError(t, json.NewDecoder(f).Decode(&got))
	assert.Equal(t, testEventFO, got)
}

func TestFileObserver_AppendsMultipleEvents(t *testing.T) {
	path := filepath.Join(t.TempDir(), "audit.log")
	obs, err := audit.NewFileObserver(path)
	require.NoError(t, err)
	defer obs.Close()

	require.NoError(t, obs.Process(testEventFO))
	require.NoError(t, obs.Process(testEventFO))

	f, err := os.Open(path)
	require.NoError(t, err)
	defer f.Close()

	var lines int
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		lines++
	}
	assert.Equal(t, 2, lines)
}

func TestFileObserver_InvalidPath(t *testing.T) {
	_, err := audit.NewFileObserver("/nonexistent/dir/audit.log")
	assert.Error(t, err)
}

func TestFileObserver_Close(t *testing.T) {
	path := filepath.Join(t.TempDir(), "audit.log")
	obs, err := audit.NewFileObserver(path)
	require.NoError(t, err)
	assert.NoError(t, obs.Close())
}
