package postgres

import (
	"context"
	"time"

	"github.com/Gunga-D/service-godzilla-soft-site/internal/item"
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

func (r *repo) CreateItem(ctx context.Context, i item.Item) (int64, error) {
	q := sq.Insert("public.item").
		Columns(
			"title",
			"description",
			"category_id",
			"platform",
			"region",
			"current_price",
			"is_for_sale",
			"old_price",
			"thumbnail_url",
			"status",
			"created_at",
			"updated_at",
		).Values(
		i.Title,
		i.Description,
		i.CategoryID,
		i.Platform,
		i.Region,
		i.CurrentPrice,
		i.IsForSale,
		i.OldPrice,
		i.ThumbnailURL,
		i.Status,
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

func (r *repo) PauseItem(ctx context.Context, itemID int64) error {
	query, args, err := sq.Update("public.item").
		Where(sq.Eq{"id": itemID}).
		Set("status", item.PausedStatus).
		Set("updated_at", time.Now()).
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

func (r *repo) FetchItemsPaginatedCursorItemId(ctx context.Context, limit uint64, cursor int64) ([]item.Item, error) {
	query, args, err := sq.Select("*").From(`public.item`).
		Where(sq.And{
			sq.Gt{"id": cursor},
			sq.Eq{"status": item.ActiveStatus},
		}).
		OrderBy("id").
		Limit(limit).
		PlaceholderFormat(sq.Dollar).ToSql()

	if err != nil {
		return nil, err
	}

	var res []item.Item
	err = r.db.SelectContext(ctx, &res, query, args...)
	if err != nil {
		return nil, err
	}

	return res, nil
}
