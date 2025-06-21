package postgres

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/Gunga-D/service-godzilla-soft-site/internal/user"
	"github.com/Gunga-D/service-godzilla-soft-site/pkg/postgres"
	"github.com/Gunga-D/service-godzilla-soft-site/pkg/redis"
	sq "github.com/Masterminds/squirrel"
)

type repo struct {
	db    postgres.TxDatabase
	redis redis.Redis
}

func NewRepo(db postgres.TxDatabase, redis redis.Redis) *repo {
	return &repo{
		db:    db,
		redis: redis,
	}
}

func (r *repo) CreateUser(ctx context.Context, usr user.User) (int64, error) {
	q := sq.Insert("public.user").
		Columns(
			"email",
			"password",
			"photo_url",
			"username",
			"first_name",
			"telegram_id",
			"steam_link",
			"has_registration_gift",
			"created_at",
			"updated_at",
		).Values(
		usr.Email,
		usr.Password,
		usr.PhotoURL,
		usr.Username,
		usr.FirstName,
		usr.TelegramID,
		usr.SteamLink,
		usr.HasRegistrationGift,
		time.Now(),
		time.Now(),
	)
	query, args, err := q.
		Suffix(`
			RETURNING id
		`).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return 0, err
	}

	var id int64
	if err := r.db.GetContext(ctx, &id, query, args...); err != nil {
		return 0, err
	}

	usr.ID = id
	userRaw, err := json.Marshal(usr)
	if err != nil {
		return 0, err
	}
	err = r.redis.Set(ctx, fmt.Sprintf(user.UserCacheKey, id), userRaw, nil)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (r *repo) GetUserByID(ctx context.Context, id int64) (*user.User, error) {
	query, args, err := sq.Select("*").From(`public.user`).
		Where(sq.Eq{"id": id}).
		Limit(1).
		PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return nil, err
	}

	var res []user.User
	if err := r.db.SelectContext(ctx, &res, query, args...); err != nil {
		return nil, err
	}
	if len(res) == 0 {
		return nil, nil
	}
	return &res[0], nil
}

func (r *repo) GetUserByEmail(ctx context.Context, email string) (*user.User, error) {
	query, args, err := sq.Select("*").From(`public.user`).
		Where(sq.Eq{"email": email}).
		Limit(1).
		PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return nil, err
	}

	var res []user.User
	if err := r.db.SelectContext(ctx, &res, query, args...); err != nil {
		return nil, err
	}
	if len(res) == 0 {
		return nil, nil
	}
	return &res[0], nil
}

func (r *repo) GetUserByTelegramID(ctx context.Context, telegramID int64) (*user.User, error) {
	query, args, err := sq.Select("*").From(`public.user`).
		Where(sq.Eq{"telegram_id": telegramID}).
		Limit(1).
		PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return nil, err
	}

	var res []user.User
	if err := r.db.SelectContext(ctx, &res, query, args...); err != nil {
		return nil, err
	}
	if len(res) == 0 {
		return nil, nil
	}
	return &res[0], nil
}

func (r *repo) RemoveFreeGift(ctx context.Context, telegramID int64) error {
	q := sq.Update("public.user").
		Where(sq.Eq{"telegram_id": telegramID}).
		Set("has_registration_gift", false)

	query, args, err := q.
		PlaceholderFormat(sq.Dollar).
		ToSql()

	if err != nil {
		return err
	}

	_, err = r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}
	return nil
}

func (r *repo) AssignSteamLinkToUser(ctx context.Context, userID int64, steamLink string) error {
	q := sq.Update("public.user").
		Where(sq.Eq{"id": userID}).
		Set("steam_link", steamLink)

	query, args, err := q.
		Suffix(`
			RETURNING *
		`).
		PlaceholderFormat(sq.Dollar).
		ToSql()

	if err != nil {
		return err
	}

	var usr user.User
	err = r.db.GetContext(ctx, &usr, query, args...)
	if err != nil {
		return err
	}

	userRaw, err := json.Marshal(usr)
	if err != nil {
		return err
	}
	err = r.redis.Set(ctx, fmt.Sprintf(user.UserCacheKey, userID), userRaw, nil)
	if err != nil {
		return err
	}
	return err
}

func (r *repo) AssignEmailToUser(ctx context.Context, userID int64, email string) error {
	q := sq.Update("public.user").
		Where(sq.Eq{"id": userID}).
		Set("email", email)

	query, args, err := q.
		Suffix(`
			RETURNING *
		`).
		PlaceholderFormat(sq.Dollar).
		ToSql()

	if err != nil {
		return err
	}

	var usr user.User
	err = r.db.GetContext(ctx, &usr, query, args...)
	if err != nil {
		return err
	}

	userRaw, err := json.Marshal(usr)
	if err != nil {
		return err
	}
	err = r.redis.Set(ctx, fmt.Sprintf(user.UserCacheKey, userID), userRaw, nil)
	if err != nil {
		return err
	}
	return err
}

func (r *repo) ChangePassword(ctx context.Context, userID int64, password string) error {
	q := sq.Update("public.user").
		Where(sq.Eq{"id": userID}).
		Set("password", password)

	query, args, err := q.
		PlaceholderFormat(sq.Dollar).
		ToSql()

	if err != nil {
		return err
	}

	_, err = r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}
	return nil
}
