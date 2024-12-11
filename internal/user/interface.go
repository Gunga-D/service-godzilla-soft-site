package user

import (
	"context"
)

type Repository interface {
	CreateUser(ctx context.Context, usr User) (int64, error)
	GetUserByEmail(ctx context.Context, email string) (*User, error)
}
