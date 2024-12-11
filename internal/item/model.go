package item

import "time"

const (
	ActiveStatus = "active"
	PausedStatus = "paused"
)

type Item struct {
	ID           int64     `db:"id"`
	Title        string    `db:"title"`
	Description  *string   `db:"description"`
	CategoryID   int64     `db:"category_id"`
	Platform     string    `db:"platform"`
	Region       string    `db:"region"`
	CurrentPrice float64   `db:"current_price"`
	IsForSale    bool      `db:"is_for_sale"`
	OldPrice     *float64  `db:"old_price"`
	ThumbnailURL string    `db:"thumbnail_url"`
	Status       string    `db:"status"`
	Slip         string    `db:"slip"`
	CreatedAt    time.Time `db:"created_at"`
	UpdatedAt    time.Time `db:"updated_at"`
}
