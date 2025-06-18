package domain

import "time"

type Request struct {
	ID        int
	Method    string //GET, POST, etc
	URL       string
	Body      []byte
	Headers   map[string]string
	CreatedAt time.Time
}
