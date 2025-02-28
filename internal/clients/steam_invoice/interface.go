package steam_invoice

import "context"

type Client interface {
	CreateInvoice(ctx context.Context, login string, amount int64) error
}
