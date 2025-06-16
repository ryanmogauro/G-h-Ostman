package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/ryanmogauro/ghostman/internal/domain"
	"github.com/ryanmogauro/ghostman/internal/infra/httpclient"
)

func main() {
	flag.Parse()
	if len(flag.Args()) > 0 {
		target := flag.Arg(0)
		req := domain.Request{Method: "GET", URL: target}

		client := httpclient.New()
		resp, err := client.Do(req)

		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("%d %s\n", resp.Status, target)
		return
	}
	fmt.Printf("Nothing to see here")
}
