package storage

import (
	"context"
	"github.com/go-redis/redis/v8"
	"log"
	"lyag/core"
	"strconv"
	"time"
)

var ctx = context.Background()

type Checker interface {
	WasSent(s *core.Seance) (bool, error)
	SetStatus(s *core.Seance, sent bool, ttl time.Duration) error
}

type redisClient struct {
	c *redis.Client
}

func (client *redisClient) SetStatus(s *core.Seance, sent bool, ttl time.Duration) error {
	key := cacheKey(s)
	log.Printf("Set %v for %s", sent, key)
	return client.c.Set(ctx, key, strconv.FormatBool(sent), ttl).Err()
}

func (client *redisClient) WasSent(s *core.Seance) (bool, error) {
	key := cacheKey(s)
	log.Printf("Checking %s", key)
	val, err := client.c.Get(ctx, key).Result()
	log.Printf("Got `%s` for %s", val, key)

	if err == redis.Nil {
		return false, nil
	} else if err != nil {
		return false, err
	}

	return strconv.ParseBool(val)
}

func NewRedis() Checker {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "redis:6379",
		Password: "", // no password set
		DB:       3,  // use default DB
	})

	cli := &redisClient{
		c: rdb,
	}
	return cli
}

func cacheKey(s *core.Seance) string {
	return "sent_" + s.GetKey()
}
