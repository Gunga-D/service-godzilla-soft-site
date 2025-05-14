package postgres

import (
	"context"
	"time"

	"github.com/Gunga-D/service-godzilla-soft-site/internal/neuro"
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

func (r *repo) CreateFinishedNeuroTask(ctx context.Context, finishedTask neuro.Task) (int64, error) {
	q := sq.Insert("public.finished_neuro_task").
		Columns(
			"id",
			"query",
			"result",
			"created_at",
		).Values(
		finishedTask.ID,
		finishedTask.Query,
		finishedTask.Result,
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
