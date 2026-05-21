package audit_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/FeshLig/metrcollector/internal/audit"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var testEventO = audit.Event{
	Timestamp: 1700000000,
	Metrics:   []string{"cpu", "mem"},
	IPAddress: "127.0.0.1",
}

func TestHTTPObserver_SendsEvent(t *testing.T) {
	var received audit.Event

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, "application/json", r.Header.Get("Content-Type"))
		require.NoError(t, json.NewDecoder(r.Body).Decode(&received))
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	obs := audit.NewHTTPObserver(srv.URL)
	require.NoError(t, obs.Process(testEventO))
	assert.Equal(t, testEventO, received)
}

func TestHTTPObserver_ServerError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer srv.Close()

	obs := audit.NewHTTPObserver(srv.URL)
	assert.Error(t, obs.Process(testEventO))
}

func TestHTTPObserver_BadURL(t *testing.T) {
	obs := audit.NewHTTPObserver("http://127.0.0.1:0")
	assert.Error(t, obs.Process(testEventO))
}
