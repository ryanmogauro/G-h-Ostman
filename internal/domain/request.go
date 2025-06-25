package domain

import (
	"flag"
	"fmt"
	"os"
	"slices"
	"strings"
	"time"
)

type headerList []string

func (h *headerList) Set(value string) error {
	*h = append(*h, value)
	return nil
}

func (h *headerList) String() string {
	return strings.Join(*h, ",")
}

type Request struct {
	ID        int
	Method    string //GET, POST, etc
	URL       string
	Body      []byte
	Headers   map[string]string
	CreatedAt time.Time
}

func ArgsToRequest(osArgs []string) (Request, error) {
	verb := strings.ToUpper(os.Args[2])
	url := os.Args[3]

	var allowedVerbs = []string{"GET", "POST", "PUT", "DELETE", "PATCH", "HEAD", "OPTIONS"}

	if !slices.Contains(allowedVerbs, verb) {
		return Request{}, fmt.Errorf("invalid verb: %s \nAllowed verbs: %v", verb, allowedVerbs)
	}

	if !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") {
		url = "https://" + url
	}

	fs := flag.NewFlagSet(verb, flag.ExitOnError)
	var data string
	var headers headerList
	var timeout time.Duration

	fs.StringVar(&data, "d", "", "Request body; prefix with @ to read from file")
	fs.Var(&headers, "H", "Custom headers")
	fs.DurationVar(&timeout, "timeout", 30*time.Second, "Request timeout")

	_ = fs.Parse(os.Args[3:])

	var body []byte
	if strings.HasPrefix(data, "@") {
		filepath := strings.TrimPrefix(data, "@")
		content, err := os.ReadFile(filepath)
		if err != nil {
			fmt.Printf("Error reading file %s: %v\n", filepath, err)
			os.Exit(1)
		}
		body = content
	} else {
		body = []byte(data)
	}

	headerMap := make(map[string]string, len(headers))
	for _, h := range headers {
		parts := strings.SplitN(h, ":", 2)
		if len(parts) == 2 {
			headerMap[strings.TrimSpace(parts[0])] = strings.TrimSpace(parts[1])
		}
	}
	req := Request{
		Method:    verb,
		URL:       url,
		Headers:   headerMap,
		Body:      body,
		CreatedAt: time.Now(),
	}

	return req, nil
}
