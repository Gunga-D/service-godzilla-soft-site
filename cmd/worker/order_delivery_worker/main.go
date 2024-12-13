package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"text/template"

	"github.com/Gunga-D/service-godzilla-soft-site/internal/clients/yandex_mail"
	order_delivery "github.com/Gunga-D/service-godzilla-soft-site/internal/order/delivery"
	order_postgres "github.com/Gunga-D/service-godzilla-soft-site/internal/order/postgres"
	"github.com/Gunga-D/service-godzilla-soft-site/pkg/postgres"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	postgres := postgres.Get(ctx)

	orderRepo := order_postgres.NewRepo(postgres)

	yandexMailClient := yandex_mail.NewClient(os.Getenv("YANDEX_MAIL_ADDRESS"),
		os.Getenv("YANDEX_MAIL_LOGIN"),
		os.Getenv("YANDEX_MAIL_PASSWORD"))

	dt, err := template.ParseFiles("assets/delivery-order-template.html")
	if err != nil {
		log.Fatalln("failed loading html templates")
	}

	orderDeliverySrv := order_delivery.NewService(orderRepo, yandexMailClient, dt)

	log.Println("start order delivery worker")
	err = orderDeliverySrv.StartOrderDelivery(ctx)
	if err != nil {
		log.Printf("[error] order delivery worker finished with an error: %v", err)
		return
	}
}
