package client

import (
	"errors"
	"os"

	"github.com/go-redis/redis"
)

const REDIS_CLIENT_ENV = "REDIS_URL"

var NoURLErr = errors.New("No URL found in " + REDIS_CLIENT_ENV)

func GetRedisFromEnv() (*redis.Client, error) {
	url := os.Getenv(REDIS_CLIENT_ENV)
	if url == "" {
		return nil, NoURLErr
	}

	return GetRedisFromUrl(url)
}

func GetRedisFromUrl(url string) (*redis.Client, error) {
	opts, err := redis.ParseURL(url)
	if err != nil {
		return nil, err
	}

	rdb := redis.NewClient(opts)
	return rdb, nil
}
