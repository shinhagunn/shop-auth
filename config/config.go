package config

func InitializeConfig() {
	InitMongoDB()
	InitRedis()
	InitSessionStore()
}
