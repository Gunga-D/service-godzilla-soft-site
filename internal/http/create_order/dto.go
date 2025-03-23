package create_order

type CreateOrderRequest struct {
	Email        *string `json:"email"`
	SteamProfile *string `json:"steam_profile"`
	Voucher      *string `json:"voucher"`
	ItemID       int64   `json:"item_id"`
}

type CreateOrderResponsePayload struct {
	OrderID     string `json:"order_id"`
	PaymentLink string `json:"payment_link"`
}
