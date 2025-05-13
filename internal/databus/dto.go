package databus

type ChangeItemStateDTO struct {
	ItemID int64  `json:"item_id"`
	Status string `json:"status"`
}

type NewUserEmailDTO struct {
	UserID *int64 `json:"user_id"`
	Email  string `json:"email"`
}

type SendToEmailDTO struct {
	Email   string `json:"email"`
	Subject string `json:"subject"`
	Body    string `json:"body"`
}

type NewUserSteamLinkDTO struct {
	UserID    int64  `json:"user_id"`
	SteamLink string `json:"steam_link"`
}
