package telegram_sign_in

type TelegramSignInRequest struct {
	AuthDate  *int64  `json:"auth_date,omitempty"`
	FirstName *string `json:"first_name,omitempty"`
	LastName  *string `json:"last_name,omitempty"`
	Hash      string  `json:"hash"`
	ID        int64   `json:"id"`
	PhotoURL  *string `json:"photo_url,omitempty"`
	Username  *string `json:"username,omitempty"`
}

type TelegramSignInResponsePayload struct {
	UserID      int64  `json:"user_id"`
	AccessToken string `json:"access_token"`
}
