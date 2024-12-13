package quick_user_registration

import (
	"context"

	"github.com/Gunga-D/service-godzilla-soft-site/internal/user"
)

type userRepo interface {
	CreateUser(ctx context.Context, usr user.User) (int64, error)
	GetUserByEmail(ctx context.Context, email string) (*user.User, error)
}
