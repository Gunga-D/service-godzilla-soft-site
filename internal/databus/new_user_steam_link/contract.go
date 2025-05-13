package new_user_steam_link

import (
	"context"

	"github.com/Gunga-D/service-godzilla-soft-site/internal/user"
)

type userRepo interface {
	GetUserByID(ctx context.Context, id int64) (*user.User, error)
	AssignSteamLinkToUser(ctx context.Context, userID int64, steamLink string) error
}
