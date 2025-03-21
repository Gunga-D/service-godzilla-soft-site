package check_voucher

type CheckVoucherReq struct {
	ItemID  int64  `json:"item_id"`
	Voucher string `json:"voucher"`
}

type CheckVoucherResp struct {
	OldPrice float64 `json:"old_price"`
	NewPrice float64 `json:"new_price"`
	Currency string  `json:"currency"`
	Warning  *string `json:"warning,omitempty"`
}
