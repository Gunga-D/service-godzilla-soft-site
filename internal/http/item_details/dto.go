package item_details

type ItemDTO struct {
	ID            int64    `json:"id"`
	Title         string   `json:"title"`
	Description   *string  `json:"description"`
	CategoryID    int64    `json:"category_id"`
	Platform      string   `json:"platform"`
	Region        string   `json:"region"`
	Publisher     *string  `json:"publisher,omitempty"`
	Creator       *string  `json:"creator,omitempty"`
	ReleaseDate   *string  `json:"release_date,omitempty"`
	CurrentPrice  float64  `json:"current_price"`
	IsForSale     bool     `json:"is_for_sale"`
	OldPrice      *float64 `json:"old_price"`
	ThumbnailURL  string   `json:"thumbnail_url"`
	BackgroundURL *string  `json:"background_url,omitempty"`
	Slip          string   `json:"slip"`
}
