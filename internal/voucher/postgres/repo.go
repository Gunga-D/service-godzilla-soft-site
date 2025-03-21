package postgres

import (
	"context"
	"time"

	"github.com/Gunga-D/service-godzilla-soft-site/internal/voucher"
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

func (r *repo) CreateVoucher(ctx context.Context, v voucher.Voucher) (int64, error) {
	q := sq.Insert("public.voucher").
		Columns(
			"type",
			"value",
			"impact",
			"has_activated",
			"created_at",
			"updated_at",
		).Values(
		v.Type,
		v.Value,
		v.Impact,
		v.HasActivated,
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

func (r *repo) ApplyVoucher(ctx context.Context, value string) (*voucher.Voucher, error) {
	v, err := r.GetActiveVoucherByValue(ctx, value)
	if err != nil {
		return nil, err
	}
	if v == nil {
		return nil, voucher.ErrNotFoundVoucher
	}
	if err := r.activateVoucher(ctx, v.ID); err != nil {
		return nil, err
	}
	return v, nil
}

func (r *repo) GetActiveVoucherByValue(ctx context.Context, value string) (*voucher.Voucher, error) {
	query, args, err := sq.Select("*").From(`public.voucher`).
		Where(sq.And{
			sq.Eq{"value": value},
			sq.Eq{"has_activated": false},
		}).
		Limit(1).
		PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return nil, err
	}

	var res []voucher.Voucher
	if err := r.db.SelectContext(ctx, &res, query, args...); err != nil {
		return nil, err
	}
	if len(res) == 0 {
		return nil, nil
	}
	return &res[0], nil
}

func (r *repo) activateVoucher(ctx context.Context, id int64) error {
	query, args, err := sq.Update("public.voucher").
		Where(sq.Eq{"id": id}).
		Set("has_activated", true).
		Set("updated_at", time.Now()).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return err
	}

	if _, err := r.db.Do(ctx).ExecContext(ctx, query, args...); err != nil {
		return err
	}
	return nil
}
