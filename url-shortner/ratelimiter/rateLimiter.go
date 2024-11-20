package ratelimiter

import (
	"context"
	"fmt"
	"net/http"
	"time"
	"urlshortner/helper"

	"github.com/go-redis/redis/v8"
)

var redisClient *redis.Client

func InitRedis() {
	redisClient = redis.NewClient(&redis.Options{
		Addr: "redis:6379", // Change this to your Redis server address
	})
	_, err := redisClient.Ping(context.Background()).Result()
	if err != nil {
		fmt.Println("Error connecting to Redis:", err)
		panic(err)
	}
	fmt.Println("Successfully connected to RateLimiter Redis")
}

func RateLimiting(r *http.Request) (bool, error) {
	fmt.Println("Rate Limiting")
	clientIP := helper.GetClientIP(r)
	fmt.Println(clientIP)
	fmt.Println("Client IP:", clientIP)

	count, err := redisClient.Incr(context.Background(), clientIP).Result()

	if err != nil {
		fmt.Println("Error incrementing key")
		return false, err
	}

	if count == 1 {
		redisClient.Expire(context.Background(), clientIP, 24*time.Hour)
	}

	if count > int64(20) {
		return false, nil // Deny request
	}
	fmt.Println("Count:", count)

	return true, nil // Allow request
}
