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

func (r *repo) CreateCode(ctx context.Context, itemID int64, value string) error {
	query, args, err := sq.Insert("public.code").
		Columns(
			"value",
			"item_id",
			"status",
			"created_at",
			"updated_at",
		).Values(
		itemID,
		value,
		code.FreeStatus,
		time.Now(),
		time.Now(),
	).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return err
	}

	if _, err := r.db.ExecContext(ctx, query, args...); err != nil {
		return err
	}
	return nil
}

func (r *repo) HasActiveCode(ctx context.Context, itemID int64) (bool, error) {
	query, args, err := sq.Select(`value`).
		From("public.code").
		Where(sq.Eq{
			"item_id": itemID,
			"status":  code.FreeStatus,
		}).
		PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return false, err
	}

	var value string
	err = r.db.GetContext(ctx, &value, query, args...)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}
