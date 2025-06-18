package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/ryanmogauro/ghostman/internal/infra/httpclient"
	"github.com/ryanmogauro/ghostman/internal/infra/storeage"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: ghostman <method> <url> [flags]")
		fmt.Println("OR ghostman history")
		os.Exit(1)
	}

	db, err := storeage.InitDB("ghostman.db")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	verb := strings.ToUpper(os.Args[2])

	switch verb {
	case "GET":

		// Add https:// if no scheme is specified
		client := httpclient.New()
		response, err := client.Send(os.Args, db)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Response: %v\n", response)

	case "HISTORY":
		requests, err := storeage.GetHistory(db)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}
		for _, request := range requests {
			fmt.Printf("ID: %v \t URL: %v \t Method: %v \t \n", request.ID, request.URL, request.Method)
		}
	}
}
