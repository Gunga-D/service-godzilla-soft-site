package user_register

type UserRegisterRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserRegisterResponsePayload struct {
	UserID      int64  `json:"user_id"`
	AccessToken string `json:"access_token"`
}
