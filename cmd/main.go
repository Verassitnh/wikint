package main

import (
	"fmt"
	"os"

	"github.com/verassitnh/wikint/cmd/internal/dp"
)

func main() {
	args := os.Args[1:]

	if len(args) < 1 {
		fmt.Println("Usage: dp - starts data pipeline; api - starts api")
		os.Exit(1)
	}

	switch args[0] {
	case "dp":
		dp.ScrapeUser(dp.StartingURL)
	case "api":
		fmt.Print("api not implemented yet")
	default:
		fmt.Println("Unknown command:", args[0])
		os.Exit(1)
	}
}
