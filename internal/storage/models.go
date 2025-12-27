package storage

import (
	"encoding/json"

	"github.com/YashIIT0909/TRexT/internal/http"
)

// Collection represents a group of requests
type Collection struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

// SavedRequest represents a request stored in the database
type SavedRequest struct {
	ID           int64  `json:"id"`
	Name         string `json:"name"`
	URL          string `json:"url"`
	Method       string `json:"method"`
	Headers      string `json:"headers"` // JSON-encoded headers
	Body         string `json:"body"`
	CollectionID int64  `json:"collection_id"`
}

// ToHTTPRequest converts a SavedRequest to an http.Request
func (sr *SavedRequest) ToHTTPRequest() *http.Request {
	headers := make(map[string]string)
	if sr.Headers != "" {
		_ = json.Unmarshal([]byte(sr.Headers), &headers)
	}

	return &http.Request{
		ID:      sr.ID,
		Name:    sr.Name,
		Method:  sr.Method,
		URL:     sr.URL,
		Headers: headers,
		Body:    sr.Body,
	}
}

// FromHTTPRequest creates a SavedRequest from an http.Request
func FromHTTPRequest(req *http.Request, collectionID int64) *SavedRequest {
	headersJSON, _ := json.Marshal(req.Headers)

	return &SavedRequest{
		ID:           req.ID,
		Name:         req.Name,
		URL:          req.URL,
		Method:       req.Method,
		Headers:      string(headersJSON),
		Body:         req.Body,
		CollectionID: collectionID,
	}
}

// HistoryEntry represents a request in history
type HistoryEntry struct {
	ID         int64  `json:"id"`
	URL        string `json:"url"`
	Method     string `json:"method"`
	StatusCode int    `json:"status_code"`
	Duration   int64  `json:"duration_ms"`
	Timestamp  int64  `json:"timestamp"`
}
