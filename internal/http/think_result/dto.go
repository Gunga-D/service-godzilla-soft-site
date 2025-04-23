package think_result

type ThinkResultRequest struct {
	ID string `json:"id"`
}

type ThinkResultResponse struct {
	Reflection string    `json:"reflection"`
	Items      []ItemDTO `json:"items"`
}

type ItemDTO struct {
	ID              int64    `json:"id"`
	CategoryID      int64    `json:"category_id"`
	Title           string   `json:"title"`
	CurrentPrice    float64  `json:"current_price"`
	ThumbnailURL    string   `json:"thumbnail_url"`
	Type            string   `json:"type"`
	HorizontalImage *string  `json:"horizontal_image,omitempty"`
	Genres          []string `json:"genres,omitempty"`
	ReleaseDate     *string  `json:"release_date,omitempty"`
}
