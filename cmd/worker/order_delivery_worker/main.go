package main

import (
	"context"
	"html/template"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/Gunga-D/service-godzilla-soft-site/internal/clients/steam_invoice"
	"github.com/Gunga-D/service-godzilla-soft-site/internal/clients/yandex_mail"
	order_delivery "github.com/Gunga-D/service-godzilla-soft-site/internal/order/delivery"
	order_postgres "github.com/Gunga-D/service-godzilla-soft-site/internal/order/postgres"
	"github.com/Gunga-D/service-godzilla-soft-site/pkg/logger"
	"github.com/Gunga-D/service-godzilla-soft-site/pkg/postgres"
	tele "gopkg.in/telebot.v4"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	telebot, err := tele.NewBot(tele.Settings{
		Token:  os.Getenv("GIFTS_TELEGRAM_TOKEN"),
		Poller: &tele.LongPoller{Timeout: 10 * time.Second},
	})
	if err != nil {
		log.Printf("[error] cant init telegram bot: %v", err)
		return
	}
	logger.Get().SetBot(telebot)

	postgres := postgres.Get(ctx)

	orderRepo := order_postgres.NewRepo(postgres)

	yandexMailClient := yandex_mail.NewClient(os.Getenv("YANDEX_MAIL_ADDRESS"),
		os.Getenv("YANDEX_MAIL_LOGIN"),
		os.Getenv("YANDEX_MAIL_PASSWORD"))

	steamInvoiceClient := steam_invoice.NewClient(os.Getenv("STEAM_INVOICE_URL"), os.Getenv("STEAM_INVOICE_TOKEN"))

	funcMap := template.FuncMap{
		"nl2br": func(text string) template.HTML {
			return template.HTML(strings.ReplaceAll(text, "\n", "<br>"))
		},
	}
	dt, err := template.New("delivery-order-template.html").Funcs(funcMap).ParseFiles("assets/delivery-order-template.html")
	if err != nil {
		log.Fatalln("failed loading html templates")
	}

	orderDeliverySrv := order_delivery.NewService(orderRepo, steamInvoiceClient, yandexMailClient, dt)

	log.Println("start order delivery worker")
	err = orderDeliverySrv.StartOrderDelivery(ctx)
	if err != nil {
		log.Printf("[error] order delivery worker finished with an error: %v", err)
		return
	}
}
