package reviews

type ReviewsResponse struct {
	Score   *float64    `json:"score,omitempty"`
	Reviews []ReviewDTO `json:"reviews,omitempty"`
}

type ReviewDTO struct {
	Comment   *string `json:"comment"`
	Score     int     `json:"score"`
	CreatedAt string  `json:"created_at"`
}
