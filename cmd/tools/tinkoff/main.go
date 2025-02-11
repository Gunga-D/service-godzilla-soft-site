package main

import (
	"context"
	"log"
	"math/rand"
	"os"

	"github.com/Gunga-D/service-godzilla-soft-site/internal/clients/tinkoff"
)

var letterRunes = []rune("0123456789")

func main() {
	c := tinkoff.NewClient(os.Getenv("TINKOFF_API_URL"), os.Getenv("TINKOFF_PASSWORD"), os.Getenv("TINKOFF_TERMINAL_KEY"))

	orderID := randStringRunes(5)
	resp, err := c.CreateInvoice(context.Background(), orderID, 5000, "Подарочная карта на 1000 рублей")
	if err != nil {
		log.Fatalf("error to create tinkoff invoice: %v", err)
	}
	log.Printf("Create invoice\nPayment link - %s\nSuccess - %v\nError Code - %s\n", resp.PaymentURL, resp.Success, resp.ErrorCode)
}

func randStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}
