package config

import (
	"context"
	"crypto/tls"
	"fmt"
	"log"
	"os"

	"github.com/redis/go-redis/v9"
)

var RedisClient *redis.Client

func InitRedis() {
	host := os.Getenv("REDIS_HOST")
	port := os.Getenv("REDIS_PORT")
	password := os.Getenv("REDIS_PASSWORD")
	if host == "" {
		host = "localhost"
	}
	if port == "" {
		port = "6379"
	}

	client := redis.NewClient(&redis.Options{
		Addr:      fmt.Sprintf("%s:%s", host, port),
		Password:  password,
		DB:        0,
		TLSConfig: &tls.Config{},
	})

	// Test connection
	_, err := client.Ping(context.Background()).Result()
	if err != nil {
		log.Println("Redis connection failed:", err)
		return // jangan crash
	}
	RedisClient = client
	log.Println("Connected to Redis successfully")
}
