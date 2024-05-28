package grafana

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

// HTTPClient is satisfied by httpClient
type HTTPClient interface {
	Get(ctx context.Context, p string, params url.Values) ([]byte, error)
}

// httpClient implements the HTTPClient interface
type httpClient struct {
	baseURL *url.URL
	client  *http.Client
	headers map[string]string
}

// Send an HTTP GET Request for the given URL path and Query parameters
func (c *httpClient) Get(ctx context.Context, p string, params url.Values) ([]byte, error) {
	// Construct the full Request URL path
	u, err := url.JoinPath(c.baseURL.String(), p)
	if err != nil {
		return nil, err
	}

	// Create a new HTTP Request object with the given Context
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, nil)
	if err != nil {
		return nil, err
	}

	// Set HTTP Headers for the Request
	for k, v := range c.headers {
		req.Header.Set(k, v)
	}

	// Set HTTP Query parameters for the Request
	req.URL.RawQuery = params.Encode()

	// Send the HTTP Request
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Unexpected HTTP Response status code: %d", resp.StatusCode)
	}

	// Read and return the Response Body
	return io.ReadAll(resp.Body)
}

// NewClient creates a new Grafana HTTP Client
func NewClient(baseURL string, accessToken string) (*httpClient, error) {
	// Parse the raw URL string into a URL structure
	u, err := url.Parse(baseURL)
	if err != nil {
		return nil, fmt.Errorf("Error parsing Grafana URL: %w", err)
	}

	client := &httpClient{
		client:  http.DefaultClient,
		baseURL: u,
		headers: map[string]string{
			"User-Agent":    "Rapid-Release-Model Metrics CLI",
			"Accept":        "application/json",
			"Content-Type":  "application/json",
			"Authorization": fmt.Sprintf("Bearer %s", accessToken),
		},
	}

	return client, nil
}
