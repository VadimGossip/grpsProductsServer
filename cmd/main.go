package main

import "github.com/VadimGossip/grpsProductsServer/internal/app"

var configDir = "config"

func main() {
	app.Run(configDir)
}
