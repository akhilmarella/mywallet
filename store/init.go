package store

import (
	"github.com/go-redis/redis"
	"github.com/rs/zerolog/log"
)

var Client *redis.Client

func InitRedis() {
	Client = redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	_, err := Client.Ping().Result()
	if err != nil {
		log.Error().Err(err).Msg("error found in redis address")
		return
	}
}
