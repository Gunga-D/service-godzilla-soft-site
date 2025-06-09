package main

import (
	"context"
	"fmt"
	topics_redis "github.com/Gunga-D/service-godzilla-soft-site/internal/topics/redis"
	"github.com/Gunga-D/service-godzilla-soft-site/pkg/redis"
	"log"
	"os/signal"
	"syscall"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	repo := redis.Get(ctx)
	redisTopicRepo := topics_redis.NewRedisRepo(repo)

	topics, err := redisTopicRepo.FetchAllTopics(ctx)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("cache contains %d topics\n", len(topics))
}
