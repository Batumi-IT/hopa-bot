package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/cenkalti/backoff/v4"
	"github.com/redis/go-redis/v9"
)

type RedisConfig struct {
	Host     string
	Port     string
	Password string
	User     string
}

// Connect to Redis and return a Redis client.
// Wait for the connection to be established before returning.
func connectToRedis(conf RedisConfig) *redis.Client {
	bf := backoff.NewExponentialBackOff()
	bf.InitialInterval = 10 * time.Second
	bf.MaxInterval = 25 * time.Second
	bf.MaxElapsedTime = 90 * time.Second

	rdb, err := backoff.RetryWithData[*redis.Client](func() (*redis.Client, error) {
		ctx, cancel := context.WithTimeout(context.Background(), bf.InitialInterval)
		defer cancel()

		conn := redis.NewClient(&redis.Options{
			Addr:     fmt.Sprintf("%s:%s", conf.Host, conf.Port),
			Password: conf.Password,
			Username: conf.User,
			DB:       0,
		})
		_, err := conn.Ping(ctx).Result()

		if err != nil {
			log.Println("Redis not yet ready...")
			return nil, err
		}
		log.Println("Connected to Redis!")
		return conn, nil
	}, bf)

	if err != nil {
		log.Fatalln(err)
		return nil
	}

	return rdb

}
