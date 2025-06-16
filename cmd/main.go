package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/ryanmogauro/ghostman/internal/domain"
	"github.com/ryanmogauro/ghostman/internal/infra/httpclient"
)

type headerList []string

func (h *headerList) String() string {
	return strings.Join(*h, ", ")
}

func (h *headerList) Set(value string) error {
	*h = append(*h, value)
	return nil
}

func main() {

	if len(os.Args) < 3 {
		fmt.Println("Usage: ghostman <method> <url> [flags]")
		os.Exit(1)
	}

	verb := strings.ToUpper(os.Args[1])
	url := os.Args[2]

	// Add https:// if no scheme is specified
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
		Method:  verb,
		URL:     url,
		Headers: headerMap,
		Body:    body,
	}

	client := httpclient.New()
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	for k, v := range resp.Headers {
		fmt.Printf("%s: %s\n", k, v)
	}

	fmt.Printf("Status: %d\n", resp.Status)
	fmt.Printf("Elapsed: %dms\n", resp.Elapsed_ms)
	fmt.Printf("Body: %s\n", string(resp.Body))
}
