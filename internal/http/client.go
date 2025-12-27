package http

import (
	"io"
	"net/http"
	"strings"
	"time"
)

// Client wraps the HTTP client with custom configuration
type Client struct {
	httpClient *http.Client
	timeout    time.Duration
}

// NewClient creates a new HTTP client with default settings
func NewClient() *Client {
	return &Client{
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		timeout: 30 * time.Second,
	}
}

// SetTimeout sets the client timeout
func (c *Client) SetTimeout(d time.Duration) {
	c.timeout = d
	c.httpClient.Timeout = d
}

// Execute performs the HTTP request and returns the response
func (c *Client) Execute(req *Request) *Response {
	startTime := time.Now()

	// Create HTTP request
	var bodyReader io.Reader
	if req.Body != "" {
		bodyReader = strings.NewReader(req.Body)
	}

	httpReq, err := http.NewRequest(req.Method, req.URL, bodyReader)
	if err != nil {
		return &Response{
			Error:    err,
			Duration: time.Since(startTime),
		}
	}

	// Set headers
	for key, value := range req.Headers {
		httpReq.Header.Set(key, value)
	}

	// Set default Content-Type for requests with body
	if req.Body != "" && httpReq.Header.Get("Content-Type") == "" {
		httpReq.Header.Set("Content-Type", "application/json")
	}

	// Execute request
	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return &Response{
			Error:    err,
			Duration: time.Since(startTime),
		}
	}
	defer resp.Body.Close()

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return &Response{
			StatusCode: resp.StatusCode,
			Status:     resp.Status,
			Headers:    resp.Header,
			Error:      err,
			Duration:   time.Since(startTime),
		}
	}

	return &Response{
		StatusCode: resp.StatusCode,
		Status:     resp.Status,
		Headers:    resp.Header,
		Body:       body,
		Duration:   time.Since(startTime),
		Size:       int64(len(body)),
	}
}

// SupportedMethods returns the list of supported HTTP methods
func SupportedMethods() []string {
	return []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"}
}
