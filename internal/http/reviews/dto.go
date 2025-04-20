package reviews

type ReviewsResponse struct {
	Score   float64     `json:"score"`
	Reviews []ReviewDTO `json:"reviews"`
}

type ReviewDTO struct {
	Comment   *string `json:"comment"`
	Score     int     `json:"score"`
	CreatedAt string  `json:"created_at"`
}
