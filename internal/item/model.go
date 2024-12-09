package item

import "time"

const (
	ActiveStatus = "active"
	PausedStatus = "paused"
)

type Item struct {
	ID           int64     `json:"id"`
	SKU          string    `json:"sku"`
	Title        string    `json:"title"`
	Description  *string   `json:"description"`
	CategoryID   int64     `json:"category_id"`
	Platform     string    `json:"platform"`
	Region       string    `json:"region"`
	CurrentPrice float64   `json:"current_price"`
	IsForSale    bool      `json:"is_for_sale"`
	OldPrice     *float64  `json:"old_price"`
	ThumbnailURL string    `json:"thumbnail_url"`
	Status       string    `json:"status"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}
