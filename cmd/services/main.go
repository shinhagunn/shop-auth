package main

import (
	"os"

	"github.com/shinhagunn/shop-auth/config"
	"github.com/shinhagunn/shop-auth/services"
)

type Service interface {
	Process()
}

func GetService(service string) Service {
	switch service {
	case "email":
		return services.NewSendEmail()
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
