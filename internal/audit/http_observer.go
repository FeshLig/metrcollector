package audit

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

// HTTPObserver sends audit events over HTTP.
type HTTPObserver struct {
	url    string
	client *http.Client
}

// NewHTTPObserver creates new HTTP audit observer.
func NewHTTPObserver(url string) *HTTPObserver {
	return &HTTPObserver{
		url:    url,
		client: &http.Client{},
	}
}

// Process sends audit event to remote HTTP endpoint.
func (h *HTTPObserver) Process(event Event) error {
	data, err := json.Marshal(event)
	if err != nil {
		return err
	}

	resp, err := h.client.Post(
		h.url,
		"application/json",
		bytes.NewReader(data),
	)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode >= http.StatusBadRequest {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return nil
}
