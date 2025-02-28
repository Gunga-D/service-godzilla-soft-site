package tinkoff

import (
	"bytes"
	"context"
	"crypto/sha256"
	"fmt"
	"log"
	"net/http"
	"sort"

	"github.com/go-resty/resty/v2"
)

type client struct {
	rc          *resty.Client
	password    string
	terminalKey string
}

func NewClient(apiURL string, password string, terminalKey string) *client {
	if apiURL == "" || password == "" || terminalKey == "" {
		log.Fatalf("tinkoff: credentials must be non-empty")
	}

	rc := resty.New()
	rc.SetBaseURL(apiURL)
	rc.Header.Set("User-Agent", "GodzillaSoft")

	return &client{
		rc:          rc,
		password:    password,
		terminalKey: terminalKey,
	}
}

func (c *client) CreateInvoice(ctx context.Context, orderID string, amount int64, description string) (*CreateInvoiceResponse, error) {
	token := generateToken(map[string]string{
		"TerminalKey": c.terminalKey,
		"OrderId":     orderID,
		"Description": description,
		"Amount":      fmt.Sprint(amount),
		"Password":    c.password,
	})

	req := CreateInvoiceRequest{
		TerminalKey: c.terminalKey,
		Amount:      amount,
		OrderID:     orderID,
		Description: description,
		Token:       token,
	}

	resp, err := c.rc.R().
		SetContext(ctx).
		SetHeader("Content-Type", "application/json").
		SetBody(req).
		SetResult(CreateInvoiceResponse{}).
		Post("/v2/Init")
	if err != nil {
		return nil, err
	}
	if resp.StatusCode() != http.StatusOK {
		return nil, fmt.Errorf("cannot create invoice: \n%s", string(resp.Body()))
	}
	log.Printf("Tinkoff body:\n%s\n", string(resp.Body()))
	return resp.Result().(*CreateInvoiceResponse), nil
}

func generateToken(v map[string]string) string {
	keys := make([]string, 0)
	for key := range v {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	var b bytes.Buffer
	for _, key := range keys {
		b.WriteString(v[key])
	}
	sum := sha256.Sum256(b.Bytes())
	return fmt.Sprintf("%x", sum)
}
