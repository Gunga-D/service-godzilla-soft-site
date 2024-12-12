package create_order

type CreateOrderRequest struct {
	Email  *string `json:"email"`
	ItemID int64   `json:"item_id"`
}

type CreateOrderResponsePayload struct {
	OrderID     string `json:"order_id"`
	PaymentLink string `json:"payment_link"`
}
