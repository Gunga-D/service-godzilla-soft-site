package admin_warmup_items

import "context"

type itemsCache interface {
	WarmUp(ctx context.Context) error
}
