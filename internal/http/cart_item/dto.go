package cart_item

type CartItemRequest struct {
	ItemID int64 `json:"item_id"`
}

type CartItemResponsePayload struct {
	Type     string  `json:"type"`
	Price    float64 `json:"price"`
	Currency string  `json:"currency"`
}
