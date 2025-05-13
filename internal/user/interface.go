package user

import (
	"context"
)

type Repository interface {
	CreateUser(ctx context.Context, usr User) (int64, error)
	GetUserByEmail(ctx context.Context, email string) (*User, error)
	GetUserByTelegramID(ctx context.Context, telegramID int64) (*User, error)
	GetUserByID(ctx context.Context, id int64) (*User, error)
	ChangePassword(ctx context.Context, userID int64, password string) error
}
