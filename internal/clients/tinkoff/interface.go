package tinkoff

import "context"

type Client interface {
	CreateInvoice(ctx context.Context, orderID string, amount int64, description string) (*CreateInvoiceResponse, error)
	CreateRecurrent(ctx context.Context, orderID string, amount int64, description string, customerKey string, notificationURL string) (*CreateRecurrentResponse, error)
	Charge(ctx context.Context, rebillID string, paymentID string) (*ChargeResponse, error)
}
