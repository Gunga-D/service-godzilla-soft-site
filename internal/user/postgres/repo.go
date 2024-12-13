package postgres

import (
	"context"
	"time"

	"github.com/Gunga-D/service-godzilla-soft-site/internal/user"
	"github.com/Gunga-D/service-godzilla-soft-site/pkg/postgres"
	sq "github.com/Masterminds/squirrel"
)

type repo struct {
	db postgres.TxDatabase
}

func NewRepo(db postgres.TxDatabase) *repo {
	return &repo{
		db: db,
	}
}

func (r *repo) CreateUser(ctx context.Context, usr user.User) (int64, error) {
	q := sq.Insert("public.user").
		Columns(
			"email",
			"password",
			"created_at",
			"updated_at",
		).Values(
		usr.Email,
		usr.Password,
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
	return id, nil
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
