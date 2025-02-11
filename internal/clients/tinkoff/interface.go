package tinkoff

import "context"

type Client interface {
	CreateInvoice(ctx context.Context, orderID string, amount int64, description string) (*CreateInvoiceResponse, error)
}
