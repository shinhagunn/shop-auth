package main

import (
	"os"

	"github.com/shinhagunn/Shop-Watches/backend/config"
	"github.com/shinhagunn/Shop-Watches/backend/services"
)

type Service interface {
	Process()
}

func GetService(service string) Service {
	switch service {
	case "email":
		return services.NewSendEmail()
	case "deliver":
		return services.NewDeliver()
	default:
		return nil
	}
}

func main() {
	config.InitializeConfig()
	service_name := os.Args[1]
	service := GetService(service_name)

	service.Process()
}
