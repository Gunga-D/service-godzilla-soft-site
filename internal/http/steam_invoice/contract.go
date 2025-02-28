package steam_invoice

import "context"

type orderCreator interface {
	CreateSteamOrder(ctx context.Context, steamLogin string, amount int64) (string, error)
}
