package popular_items

type ItemDTO struct {
	ID           int64    `json:"id"`
	Title        string   `json:"title"`
	Platform     string   `json:"platform"`
	CategoryID   int64    `json:"category_id"`
	Region       string   `json:"region"`
	CurrentPrice float64  `json:"current_price"`
	IsForSale    bool     `json:"is_for_sale"`
	OldPrice     *float64 `json:"old_price"`
	Type         string   `json:"type"`
	ThumbnailURL string   `json:"thumbnail_url"`
}
