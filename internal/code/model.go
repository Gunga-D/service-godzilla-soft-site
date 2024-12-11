package code

import "time"

const (
	FreeStatus      = "free"
	DeliveredStatus = "delivered"
)

type Code struct {
	ID        int64     `db:"id"`
	ItemID    int64     `db:"item_id"`
	Status    string    `db:"status"`
	Code      string    `db:"code"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}
