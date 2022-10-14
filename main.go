package main

import (
	"wallet-api/internal/app"

	"github.com/yanun0323/pkg/config"
)

func main() {
	if err := config.Init("config"); err != nil {
		panic(err)
	}
	app.Run()
}
