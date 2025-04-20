package postgres

import (
	"context"
	"time"

	"github.com/Gunga-D/service-godzilla-soft-site/internal/review"
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

func (r *repo) AddReview(ctx context.Context, userID *int64, itemID int64, comment *string, score int) (int64, error) {
	q := sq.Insert("public.review").
		Columns(
			"user_id",
			"item_id",
			"comment",
			"score",
			"created_at",
			"updated_at",
		).Values(
		userID,
		itemID,
		comment,
		score,
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
	if err := r.db.Do(ctx).GetContext(ctx, &id, query, args...); err != nil {
		return 0, err
	}
	return id, nil
}

func (r *repo) GetScore(ctx context.Context, itemID int64) (float64, error) {
	query, args, err := sq.Select("avg(score)").From(`public.review`).
		Where(sq.Eq{"item_id": itemID}).
		PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return 0, err
	}

	var res float64
	err = r.db.GetContext(ctx, &res, query, args...)
	if err != nil {
		return 0, err
	}
	return res, nil
}

func (r *repo) FetchCommentReviews(ctx context.Context, itemID int64, limit uint64, offset uint64) ([]review.Review, error) {
	query, args, err := sq.Select("*").From(`public.review`).
		Where(sq.And{sq.NotEq{"comment": nil}, sq.Eq{"item_id": itemID}}).
		Limit(limit).Offset(offset).
		PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return nil, err
	}

	var res []review.Review
	err = r.db.SelectContext(ctx, &res, query, args...)
	if err != nil {
		return nil, err
	}
	return res, nil
}
