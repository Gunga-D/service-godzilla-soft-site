package mtspay

type GetRateResponse struct {
	KztTopup string `json:"kzt_topup"`
	RubTopup string `json:"rub_topup"`
	UsdTopup string `json:"usd_topup"`
}
