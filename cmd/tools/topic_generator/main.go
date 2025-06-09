package main

import (
	"context"
	"fmt"
	gen "github.com/Gunga-D/service-godzilla-soft-site/internal/neuro/topics"
	topics_cached "github.com/Gunga-D/service-godzilla-soft-site/internal/topics/cached"
	topics_postgres "github.com/Gunga-D/service-godzilla-soft-site/internal/topics/postgres"
	topics_redis "github.com/Gunga-D/service-godzilla-soft-site/internal/topics/redis"
	"github.com/Gunga-D/service-godzilla-soft-site/pkg/postgres"
	"github.com/Gunga-D/service-godzilla-soft-site/pkg/redis"
	"log"

	"github.com/cohesion-org/deepseek-go"
	"os/signal"
	"syscall"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	db := postgres.Get(ctx)
	redis := redis.Get(ctx)
	cachedRepo := topics_cached.NewRepo(topics_postgres.NewRepo(db), topics_redis.NewRedisRepo(redis))

	client := deepseek.NewClient(gen.GetApiKey(), gen.GetApiURL())

	fmt.Println("Generating themes...")
	resp, err := gen.GenerateThemes(ctx, client)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Select one of the themes...")
	for id, theme := range resp.Themes {
		fmt.Printf("Theme %d:\n Title: %s\n Content: %s\n", id, theme.Title, theme.Content)
	}

	var id int
	_, err = fmt.Scan(&id)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Generating topic with title %s...\n", resp.Themes[id].Title)
	topic, err := gen.GenerateTopic(ctx, client, resp.Themes[id])
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Publishing topic...")
	identifier, err := cachedRepo.CreateTopic(ctx, topic)
	fmt.Printf("Topic successfully created with id = %d\n", identifier)
}
