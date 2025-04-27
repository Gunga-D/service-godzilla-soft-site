package fetch_collection_items

type CollectionItemDTO struct {
	ID                 int64    `json:"id"`
	Title              string   `json:"title"`
	CategoryID         int64    `json:"category_id"`
	Platform           string   `json:"platform"`
	Region             string   `json:"region"`
	CurrentPrice       float64  `json:"current_price"`
	IsForSale          bool     `json:"is_for_sale"`
	OldPrice           *float64 `json:"old_price"`
	ThumbnailURL       string   `json:"thumbnail_url"`
	HorizontalImageURL *string  `json:"horizontal_image_url"`
	ReleaseDate        *string  `json:"release_date"`
	Genres             []string `json:"genres"`
	Type               string   `json:"type"`
}
