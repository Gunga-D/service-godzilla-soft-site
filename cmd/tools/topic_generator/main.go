package main

import (
	"context"
	"fmt"
	gen "github.com/Gunga-D/service-godzilla-soft-site/internal/neuro/topics"
	pg "github.com/Gunga-D/service-godzilla-soft-site/internal/topics/postgres"
	"github.com/Gunga-D/service-godzilla-soft-site/pkg/postgres"
	"github.com/cohesion-org/deepseek-go"
	"os/signal"
	"sync"
	"syscall"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	db := postgres.Get(ctx)
	repo := pg.NewRepo(db)

	client := deepseek.NewClient(gen.GetApiKey(), gen.GetApiURL())
	var wg sync.WaitGroup
	errChan := make(chan error, 10) // buffer for all potential errors

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			topic, err := gen.GenerateTopic(ctx, client)
			if err != nil {
				errChan <- err
				return
			}

			id, err := repo.CreateTopic(ctx, topic)
			if err != nil {
				errChan <- err
				return
			}

			fmt.Printf("Topic with id = %d is created in database\n", id)
		}()
	}

	// Wait for all goroutines to complete
	go func() {
		wg.Wait()
		close(errChan)
	}()

	// Print any errors that occurred
	for err := range errChan {
		fmt.Println("Error:", err)
	}
}
