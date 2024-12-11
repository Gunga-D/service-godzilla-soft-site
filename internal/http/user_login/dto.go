package user_login

type UserLoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserLoginResponsePayload struct {
	UserID      int64  `json:"user_id"`
	AccessToken string `json:"access_token"`
}
