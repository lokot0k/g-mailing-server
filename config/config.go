package config

import (
	"log"
	"os"
	"strconv"
)

type Config struct {
	Port           string
	DBHost         string
	DBPort         string
	DBUser         string
	DBPassword     string
	DBName         string
	JwtSecret      string
	JwtExpireHours int
	RedisAddr      string
	RedisPassword  string
	RabbitMQURL    string
}

func Load() *Config {
	// если нужен .env — раскомментировать:
	// _ = godotenv.Load()

	expireHrs, err := strconv.Atoi(os.Getenv("JWT_EXPIRE_HOURS"))
	if err != nil {
		log.Printf("warning: JWT_EXPIRE_HOURS не задано или некорректно, используется 24ч: %v", err)
		expireHrs = 24
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	return &Config{
		Port:           port,
		DBHost:         os.Getenv("DB_HOST"),
		DBPort:         os.Getenv("DB_PORT"),
		DBUser:         os.Getenv("DB_USER"),
		DBPassword:     os.Getenv("DB_PASSWORD"),
		DBName:         os.Getenv("DB_NAME"),
		JwtSecret:      os.Getenv("JWT_SECRET"),
		JwtExpireHours: expireHrs,
		RedisAddr:      os.Getenv("REDIS_ADDR"),
		RedisPassword:  os.Getenv("REDIS_PASSWORD"),
		RabbitMQURL:    os.Getenv("RABBITMQ_URL"),
	}
}
