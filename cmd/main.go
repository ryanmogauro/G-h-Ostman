package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/ryanmogauro/ghostman/internal/domain"
	"github.com/ryanmogauro/ghostman/internal/infra/httpclient"
	"github.com/ryanmogauro/ghostman/internal/infra/storeage"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: ghostman <method> <url> [-d <data>] [-H <headers>] [-timeout <timeout>]")
		fmt.Println("OR ghostman history")
		fmt.Println("OR ghostman rerun <id>")
		os.Exit(1)
	}

	db, err := storeage.InitDB("ghostman.db")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	verb := strings.ToUpper(os.Args[2])

	switch verb {

	case "HISTORY":
		requests, err := storeage.GetHistory(db)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}
		for _, request := range requests {
			fmt.Printf("ID: %v \t URL: %v \t Method: %v \t \n", request.ID, request.URL, request.Method)
		}
	case "RERUN":
		id, err := strconv.Atoi(os.Args[3])
		fmt.Println("ID: ", id)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}
		request, err := storeage.GetRequest(db, id)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}

		client := httpclient.New()
		response, err := client.Send(request, db)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Response: %v\n", response)

	default:
		client := httpclient.New()

		request, err := domain.ArgsToRequest(os.Args)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}
		response, err := client.Send(request, db)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Response: %v\n", response)
	}
}
