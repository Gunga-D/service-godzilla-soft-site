package cart_item

type CartItemRequest struct {
	ItemID int64 `json:"itemID"`
}

type CartItemResponsePayload struct {
	Price       float64 `json:"price"`
	Currency    string  `json:"currency"`
	PaymentLink string  `json:"payment_link"`
}
