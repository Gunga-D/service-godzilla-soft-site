package telegram_sign_in

type jwtService interface {
	GenerateToken(userID int64, email *string) (string, error)
}
