package http

// Request represents an HTTP request to be executed
type Request struct {
	ID      int64             `json:"id"`
	Name    string            `json:"name"`
	Method  string            `json:"method"`
	URL     string            `json:"url"`
	Headers map[string]string `json:"headers"`
	Body    string            `json:"body"`
}

// NewRequest creates a new Request with default values
func NewRequest() *Request {
	return &Request{
		Method:  "GET",
		Headers: make(map[string]string),
	}
}

// Clone creates a deep copy of the request
func (r *Request) Clone() *Request {
	headers := make(map[string]string)
	for k, v := range r.Headers {
		headers[k] = v
	}
	return &Request{
		ID:      r.ID,
		Name:    r.Name,
		Method:  r.Method,
		URL:     r.URL,
		Headers: headers,
		Body:    r.Body,
	}
}
