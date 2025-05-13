package auth

import (
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type tokenClaims struct {
	jwt.StandardClaims
	UserId    int64   `json:"user_id"`
	UserEmail *string `json:"user_email,omitempty"`
}

type jwtService struct {
	secret string
}

func NewJwtService(secret string) *jwtService {
	return &jwtService{
		secret: secret,
	}
}

func (s *jwtService) GenerateToken(userID int64, email *string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &tokenClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(12 * time.Hour).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
		UserId:    userID,
		UserEmail: email,
	})
	return token.SignedString([]byte(s.secret))
}

func (s *jwtService) ParseToken(accessToken string) (int64, *string, error) {
	token, err := jwt.ParseWithClaims(accessToken, &tokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("not valid token")
		}

		return []byte(s.secret), nil
	})
	if err != nil {
		return 0, nil, err
	}

	claims, ok := token.Claims.(*tokenClaims)
	if !ok {
		return 0, nil, fmt.Errorf("invalid type of token")
	}
	return claims.UserId, claims.UserEmail, nil
}
