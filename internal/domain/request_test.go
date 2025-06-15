package domain

import (
	"encoding/json"
	"testing"
)

func TestMarshalRoundTrip(t *testing.T) {
	r := Request{
		Method: "GET",
		URL:    "https://example.com",
		Body:   []byte("test body"),
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
	}

	// Marshal to JSON
	data, err := json.Marshal(r)
	if err != nil {
		t.Fatalf("Failed to marshal request: %v", err)
	}

	// Unmarshal back
	var clone Request
	if err := json.Unmarshal(data, &clone); err != nil {
		t.Fatalf("Failed to unmarshal request: %v", err)
	}

	// Verify all fields
	if clone.Method != r.Method {
		t.Errorf("Method mismatch: got %v, want %v", clone.Method, r.Method)
	}
	if clone.URL != r.URL {
		t.Errorf("URL mismatch: got %v, want %v", clone.URL, r.URL)
	}
	if string(clone.Body) != string(r.Body) {
		t.Errorf("Body mismatch: got %v, want %v", string(clone.Body), string(r.Body))
	}
	if len(clone.Headers) != len(r.Headers) {
		t.Errorf("Headers length mismatch: got %d, want %d", len(clone.Headers), len(r.Headers))
	}
	for k, v := range r.Headers {
		if clone.Headers[k] != v {
			t.Errorf("Header %s mismatch: got %v, want %v", k, clone.Headers[k], v)
		}
	}
}
