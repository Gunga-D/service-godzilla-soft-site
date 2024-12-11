package admin_create_item

type AdminCreateItemRequest struct {
	Title        string   `json:"title"`
	Description  *string  `json:"description"`
	Slip         string   `json:"slip"`
	CategoryID   int64    `json:"category_id"`
	Platform     string   `json:"platform"`
	Region       string   `json:"region"`
	ThumbnailURL string   `json:"thumbnail_url"`
	CurrentPrice float64  `json:"current_price"`
	IsForSale    bool     `json:"is_for_sale"`
	OldPrice     *float64 `json:"old_price"`
}

type AdminCreateItemResponsePayload struct {
	ID int64 `json:"id"`
}
