package postgres

import (
	"context"
	"time"

	"github.com/Gunga-D/service-godzilla-soft-site/internal/code"
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
	var orderID string
	err := r.db.WithTx(ctx, func(ctx context.Context) error {
		codeID, err := r.blockCode(ctx, itemID)
		if err != nil {
			return err
		}

		newOrderID, err := r.insertOrder(ctx, email, amount, codeID)
		if err != nil {
			return err
		}
		orderID = newOrderID
		return nil
	})
	if err != nil {
		return "", err
	}
	return orderID, nil
}

func (r *repo) insertOrder(ctx context.Context, email string, amount int64, codeID int64) (string, error) {
	q := sq.Insert("public.item").
		Columns(
			"email",
			"code_id",
			"amount",
			"status",
			"created_at",
			"updated_at",
		).Values(
		email,
		codeID,
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
	if err := r.db.Do(ctx).GetContext(ctx, &id, query, args...); err != nil {
		return "", err
	}
	return id, nil
}

func (r *repo) blockCode(ctx context.Context, itemID int64) (int64, error) {
	selectCodeQ := sq.
		Select("id").
		From("public.code").
		Where(sq.And{
			sq.Eq{"status": code.FreeStatus},
			sq.Eq{"item_id": itemID},
		}).
		Limit(1)

	q := sq.Update("public.code").
		Where(sq.Eq{"id": subQuery(selectCodeQ)}).
		Set("status", code.BlockedStatus).
		Set("updated_at", time.Now())

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
	if err := r.db.Do(ctx).GetContext(ctx, &id, query, args...); err != nil {
		return 0, err
	}
	return id, nil
}

func subQuery(sb sq.SelectBuilder) sq.Sqlizer {
	sql, params, _ := sb.ToSql()
	return sq.Expr("("+sql+")", params)
}
