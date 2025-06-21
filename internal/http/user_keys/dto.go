package user_keys

type UserKeyDTO struct {
	ID       string `json:"id"`
	ItemName string `json:"item_name"`
	ItemSlip string `json:"item_slip"`
	Code     string `json:"code"`
	Status   string `json:"status"`
}
