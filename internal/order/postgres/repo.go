package postgres

import (
	"context"
	"time"

	"github.com/Gunga-D/service-godzilla-soft-site/internal/order"
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

func (r *repo) CreateOrder(ctx context.Context, email string, amount int64, itemID int64) (string, error) {
	q := sq.Insert("public.item").
		Columns(
			"email",
			"item_id",
			"amount",
			"status",
			"created_at",
			"updated_at",
		).Values(
		email,
		itemID,
		amount,
		order.PendingStatus,
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
		return "", err
	}

	var id string
	if err := r.db.GetContext(ctx, &id, query, args...); err != nil {
		return "", err
	}
	return id, nil
}
