package main

import (
	"github.com/yanun0323/pkg/config"

	"wallet-api/internal/app"
)

func main() {
	if err := config.Init("config", true, "../config", "../../config"); err != nil {
		panic(err)
	}
	app.Run()
}
