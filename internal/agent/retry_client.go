package agent

import (
	"context"
	"fmt"
	"net/http"
)

type RetryClient struct {
	client *http.Client
}

func NewRetryClient(client *http.Client) *RetryClient {
	return &RetryClient{client: client}
}

func (r *RetryClient) Do(newReq func() (*http.Request, error)) (*http.Response, error) {
	var resp *http.Response

	err := withRetry(context.TODO(), func() error {
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
