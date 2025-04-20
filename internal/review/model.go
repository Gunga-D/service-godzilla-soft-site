package review

import "time"

type Review struct {
	ID        string    `db:"id"`
	UserID    *int64    `db:"user_id"`
	ItemID    int64     `db:"item_id"`
	Comment   *string   `db:"comment"`
	Score     int       `db:"score"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}
