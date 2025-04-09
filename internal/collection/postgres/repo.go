package postgres

import (
	"context"

	"github.com/Gunga-D/service-godzilla-soft-site/internal/collection"
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

func (r *repo) FetchCollectionsByFilter(ctx context.Context, criteria sq.And, limit uint64, offset uint64) ([]collection.Collection, error) {
	query, args, err := sq.Select("*").From(`public.collection`).
		Where(criteria).
		OrderBy("id").
		Limit(limit).
		Offset(offset).
		PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return nil, err
	}

	var res []collection.Collection
	err = r.db.SelectContext(ctx, &res, query, args...)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (r *repo) FetchCollectionItemsByCollectionID(ctx context.Context, collectionID int64, limit uint64, offset uint64) ([]collection.CollectionItem, error) {
	query, args, err := sq.Select("*").From(`public.collection_item`).
		Where(sq.And{
			sq.Eq{"collection_id": collectionID},
		}).
		OrderBy("id").
		Limit(limit).
		Offset(offset).
		PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return nil, err
	}

	var res []collection.CollectionItem
	err = r.db.SelectContext(ctx, &res, query, args...)
	if err != nil {
		return nil, err
	}
	return res, nil
}
