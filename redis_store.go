package main

import (
	"os"

	"github.com/fzzy/radix/redis"
)

const (
	RedisHost   = "localhost"
	RedisPort   = "6379"
	UrlStore    = "canoe:urlstore"
)

func RedisClient() *redis.Client {
	client, err := redis.Dial("tcp", redisURL())
	if err != nil {
		panic("Could not connect to redis on " + redisURL())
	}
	return client
}

func redisURL() string {
	if redisEnv := os.Getenv("REDIS_SERVER"); len(redisEnv) > 1 {
		return redisEnv
	}
	return RedisHost + ":" + RedisPort
}
