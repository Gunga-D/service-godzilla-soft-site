package postgres

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/Gunga-D/service-godzilla-soft-site/internal/code"
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

func (r *repo) CreateCode(ctx context.Context, itemID int64, value string) (int64, error) {
	q := sq.Insert("public.code").
		Columns(
			"item_id",
			"value",
			"status",
			"created_at",
			"updated_at",
		).Values(
		itemID,
		value,
		code.FreeStatus,
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

func (r *repo) HasActiveCode(ctx context.Context, itemID int64) (bool, error) {
	query, args, err := sq.Select(`id`).
		From("public.code").
		Where(sq.Eq{
			"item_id": itemID,
			"status":  code.FreeStatus,
		}).
		PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return false, err
	}

	var id int
	err = r.db.GetContext(ctx, &id, query, args...)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}
