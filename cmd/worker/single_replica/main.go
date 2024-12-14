package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"text/template"

	"github.com/Gunga-D/service-godzilla-soft-site/internal/clients/yandex_mail"
	"github.com/Gunga-D/service-godzilla-soft-site/internal/databus"
	"github.com/Gunga-D/service-godzilla-soft-site/internal/databus/change_item_state"
	"github.com/Gunga-D/service-godzilla-soft-site/internal/databus/quick_user_registration"
	"github.com/Gunga-D/service-godzilla-soft-site/internal/databus/send_to_email"
	item_postgres "github.com/Gunga-D/service-godzilla-soft-site/internal/item/postgres"
	user_postgres "github.com/Gunga-D/service-godzilla-soft-site/internal/user/postgres"
	"github.com/Gunga-D/service-godzilla-soft-site/pkg/postgres"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	postgres := postgres.Get(ctx)

	databusClient := databus.NewClient(ctx)

	itemRepo := item_postgres.NewRepo(postgres)
	userRepo := user_postgres.NewRepo(postgres)

	yandexMailClient := yandex_mail.NewClient(os.Getenv("YANDEX_MAIL_ADDRESS"),
		os.Getenv("YANDEX_MAIL_LOGIN"),
		os.Getenv("YANDEX_MAIL_PASSWORD"))

	rt, err := template.ParseFiles("assets/registration-template.html")
	if err != nil {
		log.Fatalln("failed loading html templates")
	}

	log.Println("start consume change item state databus")
	go change_item_state.NewHandler(databusClient, itemRepo).Consume(ctx)

	log.Println("start consume quick user registration databus")
	go quick_user_registration.NewHandler(databusClient, userRepo, rt, databusClient).Consume(ctx)

	log.Println("start consume send to email databus")
	go send_to_email.NewHandler(databusClient, yandexMailClient).Consume(ctx)

	<-ctx.Done()
}
