package main

import (
	"os"

	"github.com/GeoDB-Limited/geo-provider/internal/cli"
)

func main() {
	if !cli.Run(os.Args) {
		os.Exit(1)
	}
}
