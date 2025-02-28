package steam_invoice

type CreateInvoiceRequest struct {
	SteamLogin string  `json:"steam_login"`
	Amount     float64 `json:"amount"`
}
