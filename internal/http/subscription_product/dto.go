package subscription_product

type SubscriptionProductRequest struct {
	ItemID int64 `json:"item_id"`
}

type SubscriptionProductResponse struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}
