package domain

type Response struct {
	Status     int
	Headers    map[string]string
	Body       []byte
	Elapsed_ms int64
}
