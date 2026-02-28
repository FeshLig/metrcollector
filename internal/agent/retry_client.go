package agent

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"strings"
	"time"
)

type RetryClient struct {
	client *http.Client
}

func NewRetryClient(client *http.Client) *RetryClient {
	return &RetryClient{client: client}
}

func (r *RetryClient) Do(ctx context.Context, newReq func() (*http.Request, error)) (*http.Response, error) {
	var resp *http.Response

	err := withRetry(ctx, func() error {
		req, err := newReq()
		if err != nil {
			return err
		}

		tmpResp, err := r.client.Do(req)
		if err != nil {
			return err
		}

		if tmpResp.StatusCode >= 500 {
			tmpResp.Body.Close()
			return fmt.Errorf("server error: %d", tmpResp.StatusCode)
		}

		resp = tmpResp
		return nil
	})

	if err != nil {
		return nil, err
	}

	return resp, nil
}

func withRetry(ctx context.Context, fn func() error) error {
	const maxRetries = 3
	var lastErr error

	for attempt := 0; attempt <= maxRetries; attempt++ {

		err := fn()
		if err == nil {
			return nil
		}

		if !isRetriable(err) {
			return err
		}

		lastErr = err

		if attempt == maxRetries {
			break
		}

		backoff := time.Duration(1+2*attempt) * time.Second

		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(backoff):
		}
	}

	return lastErr
}

func isRetriable(err error) bool {
	var netErr net.Error
	if errors.As(err, &netErr) {
		return true
	}
	if strings.Contains(err.Error(), "server error") {
		return true
	}

	return false
}
