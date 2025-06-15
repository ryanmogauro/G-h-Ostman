package domain

type Request struct {
	Method  string //GET, POST, etc
	URL     string
	Body    []byte
	Headers map[string]string
}
