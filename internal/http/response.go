package http

import (
	"net/http"
	"time"
)

// Response represents an HTTP response
type Response struct {
	StatusCode int
	Status     string
	Headers    http.Header
	Body       []byte
	Duration   time.Duration
	Size       int64
	Error      error
}

// IsSuccess returns true if the status code is 2xx
func (r *Response) IsSuccess() bool {
	return r.StatusCode >= 200 && r.StatusCode < 300
}

// IsError returns true if the response contains an error
func (r *Response) IsError() bool {
	return r.Error != nil
}

// BodyString returns the response body as a string
func (r *Response) BodyString() string {
	return string(r.Body)
}
