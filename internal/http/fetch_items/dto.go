package fetch_items

type ItemDTO struct {
	ID           int64    `json:"id"`
	Title        string   `json:"title"`
	CategoryID   int64    `json:"category_id"`
	Platform     string   `json:"platform"`
	Region       string   `json:"region"`
	CurrentPrice float64  `json:"current_price"`
	IsForSale    bool     `json:"is_for_sale"`
	OldPrice     *float64 `json:"old_price"`
	ThumbnailURL string   `json:"thumbnail_url"`
}
