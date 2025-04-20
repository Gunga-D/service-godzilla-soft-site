package add_review

type AddReviewRequest struct {
	ItemID  int64   `json:"item_id"`
	Score   int     `json:"score"`
	Comment *string `json:"comment,omitempty"`
}

type AddReviewResponse struct {
	ReviewID int64 `json:"review_id"`
}
