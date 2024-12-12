package user_login

type jwtService interface {
	GenerateToken(userID int64, email string) (string, error)
}
