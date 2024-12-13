package databus

type ChangeItemStateDTO struct {
	ItemID int64  `json:"item_id"`
	Status string `json:"status"`
}

type QuickUserRegistrationDTO struct {
	Email string `json:"email"`
}
