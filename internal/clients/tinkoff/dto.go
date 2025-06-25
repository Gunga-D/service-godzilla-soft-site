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

type CreateRecurrentRequest struct {
	TerminalKey     string `json:"TerminalKey"`
	Amount          int64  `json:"Amount"`
	OrderID         string `json:"OrderId"`
	Token           string `json:"Token"`
	Description     string `json:"Description"`
	Recurrent       string `json:"Recurrent"`
	CustomerKey     string `json:"CustomerKey"`
	PayType         string `json:"PayType"`
	NotificationURL string `json:"NotificationURL"`
}

type CreateRecurrentResponse struct {
	Success     bool   `json:"Success"`
	ErrorCode   string `json:"ErrorCode"`
	TerminalKey string `json:"TerminalKey"`
	Status      string `json:"Status"`
	PaymentId   string `json:"PaymentId"`
	OrderId     string `json:"OrderId"`
	Amount      int64  `json:"Amount"`
	PaymentURL  string `json:"PaymentURL"`
}

type ChargeRequest struct {
	TerminalKey string `json:"TerminalKey"`
	PaymentId   string `json:"PaymentId"`
	RebillId    string `json:"RebillId"`
	Token       string `json:"Token"`
}

type ChargeResponse struct {
	TerminalKey string `json:"TerminalKey"`
	Amount      int64  `json:"Amount"`
	OrderId     string `json:"OrderId"`
	Success     bool   `json:"Success"`
	Status      string `json:"Status"`
	PaymentId   string `json:"PaymentId"`
	ErrorCode   string `json:"ErrorCode"`
}
