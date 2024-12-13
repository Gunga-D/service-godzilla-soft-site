package change_item_state

import "context"

type itemRepo interface {
	ChangeItemState(ctx context.Context, itemID int64, status string) error
}
