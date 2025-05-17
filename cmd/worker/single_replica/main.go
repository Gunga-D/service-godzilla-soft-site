package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"text/template"

	"github.com/Gunga-D/service-godzilla-soft-site/internal/clients/deepseek"
	"github.com/Gunga-D/service-godzilla-soft-site/internal/clients/steam"
	"github.com/Gunga-D/service-godzilla-soft-site/internal/clients/yandex_mail"
	"github.com/Gunga-D/service-godzilla-soft-site/internal/databus"
	"github.com/Gunga-D/service-godzilla-soft-site/internal/databus/change_item_state"
	"github.com/Gunga-D/service-godzilla-soft-site/internal/databus/neuro_new_items"
	"github.com/Gunga-D/service-godzilla-soft-site/internal/databus/new_user_email"
	"github.com/Gunga-D/service-godzilla-soft-site/internal/databus/new_user_steam_link"
	"github.com/Gunga-D/service-godzilla-soft-site/internal/databus/send_to_email"
	item_postgres "github.com/Gunga-D/service-godzilla-soft-site/internal/item/postgres"
	user_postgres "github.com/Gunga-D/service-godzilla-soft-site/internal/user/postgres"
	"github.com/Gunga-D/service-godzilla-soft-site/pkg/postgres"
	"github.com/Gunga-D/service-godzilla-soft-site/pkg/redis"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	postgres := postgres.Get(ctx)
	redis := redis.Get(ctx)

	databusClient := databus.NewClient(ctx)

	itemRepo := item_postgres.NewRepo(postgres)
	userRepo := user_postgres.NewRepo(postgres, redis)

	steamClient := steam.NewClient(os.Getenv("STEAM_KEY"))
	deepseekClient := deepseek.NewClient(os.Getenv("DEEPSEEK_URL"), os.Getenv("DEEPSEEK_TOKEN"))
	yandexMailClient := yandex_mail.NewClient(os.Getenv("YANDEX_MAIL_ADDRESS"),
		os.Getenv("YANDEX_MAIL_LOGIN"),
		os.Getenv("YANDEX_MAIL_PASSWORD"))

	rt, err := template.ParseFiles("assets/registration-template.html")
	if err != nil {
		log.Fatalln("failed loading html templates")
	}

	log.Println("start consume change item state databus")
	go change_item_state.NewHandler(databusClient, itemRepo).Consume(ctx)

	log.Println("start consume new user email databus")
	go new_user_email.NewHandler(databusClient, userRepo, rt, databusClient).Consume(ctx)

	log.Println("start consume new user email databus")
	go new_user_email.NewHandler(databusClient, userRepo, rt, databusClient).Consume(ctx)

	log.Println("start consume send to email databus")
	go send_to_email.NewHandler(databusClient, yandexMailClient).Consume(ctx)

	log.Println("start consume new user steam link databus")
	go new_user_steam_link.NewHandler(databusClient, userRepo).Consume(ctx)

	log.Println("start consume neuro new items")
	go neuro_new_items.NewHandler(itemRepo, deepseekClient, steamClient, databusClient).Consume(ctx)

	<-ctx.Done()
}
