package tinkoff

type CreateInvoiceRequest struct {
	TerminalKey string `json:"TerminalKey"`
	Amount      int64  `json:"Amount"`
	OrderID     string `json:"OrderId"`
	Token       string `json:"Token"`
	Description string `json:"Description"`
}

type CreateInvoiceResponse struct {
	Success     bool   `json:"Success"`
	ErrorCode   string `json:"ErrorCode"`
	TerminalKey string `json:"TerminalKey"`
	Status      string `json:"Status"`
	PaymentId   string `json:"PaymentId"`
	OrderId     string `json:"OrderId"`
	Amount      int64  `json:"Amount"`
	PaymentURL  string `json:"PaymentURL"`
}
