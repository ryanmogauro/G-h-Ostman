package httpclient

import (
	"bytes"
	"database/sql"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/ryanmogauro/ghostman/internal/domain"
	"github.com/ryanmogauro/ghostman/internal/infra/storeage"
)

type Client struct {
	hc *http.Client
}

type headerList []string

func New() *Client {
	return &Client{
		hc: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (c *Client) Do(req domain.Request) (domain.Response, error) {
	//build request
	method := req.Method
	url := req.URL
	body := bytes.NewReader(req.Body)
	r, err := http.NewRequest(method, url, body)

	if err != nil {
		return domain.Response{}, err
	}

	for k, v := range req.Headers {
		r.Header.Set(k, v)
	}

	//send request
	start := time.Now()
	resp, err := c.hc.Do(r)
	if err != nil {
		return domain.Response{}, err
	}

	defer resp.Body.Close()

	//read and create response
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return domain.Response{}, fmt.Errorf("failed to read response body: %w", err)
	}

	headers := make(map[string]string, len(resp.Header))
	for k, v := range resp.Header {
		headers[k] = strings.Join(v, ", ")
	}

	return domain.Response{
		Status:     resp.StatusCode,
		Headers:    headers,
		Body:       respBody,
		Elapsed_ms: time.Since(start).Milliseconds(),
	}, nil
}

func (c *Client) Send(args []string, db *sql.DB) (domain.Response, error) {
	if len(os.Args) < 4 {
		fmt.Println("Usage: ghostman <method> <url> [flags]")
		fmt.Println("Please include a target url")
	}

	verb := strings.ToUpper(os.Args[2])
	url := os.Args[3]

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
	req := domain.Request{
		Method:    verb,
		URL:       url,
		Headers:   headerMap,
		Body:      body,
		CreatedAt: time.Now(),
	}

	err := storeage.InsertRequest(db, req)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	client := New()
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	for k, v := range resp.Headers {
		fmt.Printf("%s: %s\n", k, v)
	}

	return resp, nil
}

func (h *headerList) String() string {
	return strings.Join(*h, ", ")
}

func (h *headerList) Set(value string) error {
	*h = append(*h, value)
	return nil
}
