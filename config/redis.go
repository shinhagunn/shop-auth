package config

import (
	"log"

	"github.com/shinhagunn/Shop-Watches/backend/services"
)

var RedisClient *services.RedisClient

func InitRedis() {
	RedisClient = services.NewRedisClient("localhost:6379")
	log.Println("Connected to Redis!")
}
