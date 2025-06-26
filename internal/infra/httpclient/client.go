package httpclient

import (
	"bytes"
	"database/sql"
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

func (c *Client) Send(req domain.Request, db *sql.DB) (domain.Response, error) {

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

	return resp, nil
}
