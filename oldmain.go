package main

import (
	"os"

	"github.com/GeoDB-Limited/geo-provider/app"
	"github.com/GeoDB-Limited/geo-provider/config"
)

func main() {
	cfg := config.New(os.Getenv("CONFIG"))

	if err := app.New(cfg).Run(); err != nil {
		panic(err)
	}
}
