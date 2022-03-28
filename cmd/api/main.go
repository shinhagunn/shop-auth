package main

import (
	"github.com/shinhagunn/Shop-Watches/backend/config"
	"github.com/shinhagunn/Shop-Watches/backend/routes"
)

func main() {
	config.InitializeConfig()
	routes.InitRouter()
}
