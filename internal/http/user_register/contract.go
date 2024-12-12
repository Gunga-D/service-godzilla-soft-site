package user_register

type jwtService interface {
	GenerateToken(userID int64, email string) (string, error)
}
