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

type NeuroTaskDTO struct {
	ID    string `json:"id"`
	Query string `json:"query"`
}

type NeuroNewItemsDTO struct {
	Query string `json:"query"`
}

type TelegramRegistrationDTO struct {
	TelegramID int64 `json:"telegram_id"`
}
