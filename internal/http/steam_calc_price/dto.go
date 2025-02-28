package steam_calc_price

type SteamCalcPriceRequest struct {
	Amount int64 `json:"amount"`
}

type SteamInvoiceResponse struct {
	Price int64 `json:"price"`
}
