package postgres

import (
	"context"
	"time"

	"github.com/Gunga-D/service-godzilla-soft-site/internal/subscription"
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

func (r *repo) CreateSubscriptionBill(ctx context.Context, userID int64, amount int64, forPeriod time.Duration) (string, error) {
	q := sq.Insert("public.subscription_bills").
		Columns(
			"user_id",
			"amount",
			"status",
			"created_at",
			"updated_at",
			"expired_at",
		).Values(
		userID,
		amount,
		subscription.PendingStatus,
		time.Now(),
		time.Now(),
		time.Now().Add(forPeriod).Unix(),
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

func (r *repo) PaidSubscriptionBill(ctx context.Context, id string, rebillID string) error {
	query, args, err := sq.Update("public.subscription_bills").
		Where(sq.Eq{"id": id}).
		Set("status", subscription.PaidStatus).
		Set("rebill_id", rebillID).
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

func (r *repo) GetLastUserSubscriptionBill(ctx context.Context, userID int64) (*subscription.UserSubscription, error) {
	query, args, err := sq.Select("status, expired_at").From(`public.subscription_bills`).
		Where(sq.Eq{"user_id": userID}).
		OrderBy("created_at desc").
		Limit(1).
		PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return nil, err
	}

	var res []subscription.UserSubscription
	if err := r.db.SelectContext(ctx, &res, query, args...); err != nil {
		return nil, err
	}
	if len(res) == 0 {
		return nil, nil
	}
	return &res[0], nil
}

func (r *repo) FailedSubscriptionBill(ctx context.Context, id string) error {
	query, args, err := sq.Update("public.subscription_bills").
		Where(sq.Eq{"id": id}).
		Set("status", subscription.FailedStatus).
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

func (r *repo) FetchLastUserSubscriptionBills(ctx context.Context) ([]subscription.PaidSubscription, error) {
	var res []subscription.PaidSubscription
	if err := r.db.Do(ctx).SelectContext(ctx, &res, `
		select distinct on (user_id) user_id, expired_at, rebill_id, need_prolong from public.subscription_bills
		where status=$1 order by user_id, created_at DESC
	`, subscription.PaidStatus); err != nil {
		return nil, err
	}
	return res, nil
}

func (r *repo) GetSubscriptionProduct(ctx context.Context, itemID int64) (*subscription.SubscriptionProduct, error) {
	query, args, err := sq.Select("login, password").From(`public.accounts`).
		Where(sq.Eq{"item_id": itemID}).
		Limit(1).
		PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return nil, err
	}

	var res []subscription.SubscriptionProduct
	if err := r.db.SelectContext(ctx, &res, query, args...); err != nil {
		return nil, err
	}
	if len(res) == 0 {
		return nil, nil
	}
	return &res[0], nil
}
