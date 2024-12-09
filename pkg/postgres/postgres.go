package postgres

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
)

type txKey struct{}

type Postgres struct {
	*sqlx.DB
}

func New(db *sqlx.DB) *Postgres {
	return &Postgres{
		db,
	}
}

func (p *Postgres) Do(ctx context.Context) Database {
	if tx := extractTx(ctx); tx != nil {
		return tx
	}
	return p
}

func (p *Postgres) WithTx(ctx context.Context, fn func(ctx context.Context) error) (err error) {
	if tx := extractTx(ctx); tx != nil {
		return fn(ctx)
	}
	tx, err := p.BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin tx: %v", err)
	}
	ctx = injectTx(ctx, tx)

	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback()
			panic(p)
		} else if err != nil {
			err = fmt.Errorf("transaction rolled back because of error: %w", err)
			_ = tx.Rollback()
		} else {
			err = tx.Commit()
			if err != nil {
				err = fmt.Errorf("failed to commit transaction: %w", err)
			}
		}
	}()

	err = fn(ctx)
	return err
}

func injectTx(ctx context.Context, tx *sqlx.Tx) context.Context {
	return context.WithValue(ctx, txKey{}, tx)
}

func extractTx(ctx context.Context) *sqlx.Tx {
	if tx, ok := ctx.Value(txKey{}).(*sqlx.Tx); ok {
		return tx
	}
	return nil
}
