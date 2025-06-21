package user

import "time"

const (
	UserCacheKey = "user:%d"
)

type MetaUserIDKey struct{}

type MetaUserEmailKey struct{}

type User struct {
	ID                  int64     `db:"id"`
	Email               *string   `db:"email"`
	Password            *string   `db:"password"`
	PhotoURL            *string   `db:"photo_url"`
	Username            *string   `db:"username"`
	FirstName           *string   `db:"first_name"`
	TelegramID          *int64    `db:"telegram_id"`
	SteamLink           *string   `db:"steam_link"`
	HasRegistrationGift bool      `db:"has_registration_gift"`
	CreatedAt           time.Time `db:"created_at"`
	UpdatedAt           time.Time `db:"updated_at"`
}
