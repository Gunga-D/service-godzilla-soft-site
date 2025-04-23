package think

type ThinkRequest struct {
	Query string `json:"query"`
}

type ThinkResponse struct {
	ID string `json:"id"`
}
