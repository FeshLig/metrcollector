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
	var err error

	err = withRetry(context.TODO(), func() error {

		req, err := newReq()
		if err != nil {
			return err
		}

		resp, err = r.client.Do(req)
		if err != nil {
			return err
		}

		if resp.StatusCode >= 500 {
			resp.Body.Close()
			return fmt.Errorf("server error: %d", resp.StatusCode)
		}

		return nil
	})

	return resp, err
}
