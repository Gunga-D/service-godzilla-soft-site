package steam_invoice

type SteamInvoiceRequest struct {
	SteamLogin string `json:"steam_login"`
	Amount     int64  `json:"amount"`
}

type SteamInvoiceResponse struct {
	OrderID     string `json:"order_id"`
	PaymentLink string `json:"payment_link"`
}
