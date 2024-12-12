package mdw

type jwtService interface {
	ParseToken(accessToken string) (int64, string, error)
}
