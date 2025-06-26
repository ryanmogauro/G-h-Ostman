package domain

import "fmt"

type Response struct {
	Status     int
	Headers    map[string]string
	Body       []byte
	Elapsed_ms int64
}

func (r *Response) FormatResponse() (string, error) {
	response := ""
	response += fmt.Sprintf("Status: %d\n", r.Status)

	// Only include body if it's not null/empty
	if len(r.Body) > 0 {
		bodyString := string(r.Body)
		response += fmt.Sprintf("Body: %s\n", bodyString)
	}

	response += fmt.Sprintf("Elapsed: %vms", r.Elapsed_ms)

	return response, nil
}
