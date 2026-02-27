package store

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
)

// struct wrapper around redis client
type StorageService struct {
	redisClient *redis.Client
}

var (
	storeService = &StorageService{}
	ctx          = context.Background()
)

const CacheDuration = 6 * time.Hour

// Init the store service and return store pointer
func InitializeStore() *StorageService {
	redisClient := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	pong, err := redisClient.Ping(ctx).Result()
	if err != nil {
		panic(fmt.Sprintf("Error init redis: %v", err))
	}
	fmt.Printf("\nRedis started successfully: pong message = {%s}", pong)
	storeService.redisClient = redisClient
	return storeService
}

/*
save the mapping between the original url and generated short url
*/
func SaveUrlMapping(shortUrl string, originalUrl string, UserId string) {

	err := storeService.redisClient.Set(ctx, shortUrl, originalUrl, CacheDuration).Err()
	if err != nil {
		panic(fmt.Sprintf("Failed saving key url | Error: %v - shortUrl: %s - originalUrl: %s\n", err, shortUrl, originalUrl))
	}
}

/*
We should be able to get the initial long URL once the short one is provided.
We should retrieve the long url and redirect.
*/
func RetrieveInitialUrl(shortUrl string) string {

	result, err := storeService.redisClient.Get(ctx, shortUrl).Result()

	if err != nil {
		panic(fmt.Sprintf("Failed RetrieveInitialUrl url | Error: %v - shortUrl: %s\n", err, shortUrl))
	}

	return result
}
