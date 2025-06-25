package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/Gunga-D/service-godzilla-soft-site/internal/clients/tinkoff"
	"github.com/Gunga-D/service-godzilla-soft-site/internal/subscription/checker"
	sub_repo "github.com/Gunga-D/service-godzilla-soft-site/internal/subscription/postgres"
	"github.com/Gunga-D/service-godzilla-soft-site/pkg/postgres"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	tinkoffClient := tinkoff.NewClient(os.Getenv("TINKOFF_API_URL"), os.Getenv("TINKOFF_PASSWORD"), os.Getenv("TINKOFF_TERMINAL_KEY"))

	postgres := postgres.Get(ctx)
	subRepo := sub_repo.NewRepo(postgres)

	subChecker := checker.NewService(subRepo, tinkoffClient)
	subChecker.StartCheck(ctx)
}
