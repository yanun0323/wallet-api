package main

import (
	"wallet-api/internal/app"
	"wallet-api/pkg/config"
)

func main() {
	if err := config.Init("config"); err != nil {
		panic(err)
	}

	app.Run()
}
