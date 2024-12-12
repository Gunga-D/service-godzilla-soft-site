package databus

type ItemOutOfStockDTO struct {
	ItemID int64 `json:"item_id"`
}

type QuickUserRegistrationDTO struct {
	Email string `json:"email"`
}
