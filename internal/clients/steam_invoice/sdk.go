package steam_invoice

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/go-resty/resty/v2"
)

type client struct {
	rc *resty.Client
}

func NewClient(apiURL string, token string) *client {
	if apiURL == "" || token == "" {
		log.Fatalf("steam invoice: credentials must be non-empty")
	}

	rc := resty.New()
	rc.SetBaseURL(apiURL)
	rc.Header.Set("User-Agent", "GodzillaSoft")
	rc.Header.Set("Authorization", token)

	return &client{
		rc: rc,
	}
}

func (c *client) CreateInvoice(ctx context.Context, login string, amount int64) error {
	req := CreateInvoiceRequest{
		SteamLogin: login,
		Amount:     float64(amount / 100),
	}
	resp, err := c.rc.R().
		SetContext(ctx).
		SetHeader("Content-Type", "application/json").
		SetBody(req).
		Post("/create_invoice")
	if err != nil {
		return err
	}
	if resp.StatusCode() != http.StatusOK {
		return fmt.Errorf("cannot create steam invoice, status code - %d", resp.StatusCode())
	}
	return nil
}
