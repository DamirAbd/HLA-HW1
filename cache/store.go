package cache

import (
	"context"
	"encoding/json"
	"log"

	"github.com/DamirAbd/HLA-HW1/types"
	"github.com/redis/go-redis/v9"
)

type redisCache struct {
	client *redis.Client
}

func NewRedisCache() types.FeedCache {
	client := redis.NewClient(&redis.Options{
		Addr:     "redis:6379",
		Password: "pwd",
		DB:       1,
	})

	// Check connection
	pong, err := client.Ping(context.Background()).Result()
	if err != nil {
		log.Fatalf("Can't connect to Redis: %v", err)
	}
	log.Println("Redis connected:", pong)

	return &redisCache{
		client: client,
	}
}

func (cache *redisCache) Set(key string, value []*types.Post) {

	json, err := json.Marshal(value)
	if err != nil {
		panic(err)
	}

	cache.client.Set(context.Background(), key, json, 0)
}

func (cache *redisCache) Get(key string) []*types.Post {

	val, err := cache.client.Get(context.Background(), key).Result()
	if err != nil {
		return nil
	}

	var posts []*types.Post
	err = json.Unmarshal([]byte(val), &posts)
	if err != nil {
		panic(err)
	}
	return posts
}
