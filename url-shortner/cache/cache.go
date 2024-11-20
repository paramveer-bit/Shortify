package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
	"urlshortner/model"

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
	fmt.Println("Successfully connected to Redis")
}

func Set(user model.User) error {
	fmt.Println("runnihn")
	ctx := context.Background()
	userData, err := json.Marshal(user)
	if err != nil {
		return err
	}

	// Cache the user data for 24 hours
	err = redisClient.Set(ctx, user.ShortUrl, userData, 24*time.Hour).Err()
	if err != nil {
		return err
	}
	fmt.Println("Successfully cached user data")
	return nil
}

func Get(ShortUrl string) (model.User, error) {
	ctx := context.Background()
	val, err := redisClient.Get(ctx, ShortUrl).Result()
	if err != nil {
		return model.User{}, err
	}

	var user model.User
	err = json.Unmarshal([]byte(val), &user)
	if err != nil {
		return model.User{}, err
	}

	fmt.Println("Cache Hit WonderFull")
	return user, nil
}
