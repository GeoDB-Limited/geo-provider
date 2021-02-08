package main

import (
	"github.com/geo-provider/app"
	"github.com/geo-provider/config"
	"os"
)

func main() {
	cfg := config.New(os.Getenv("CONFIG"))

	if err := app.New(cfg).Run(); err != nil {
		panic(err)
	}
}
