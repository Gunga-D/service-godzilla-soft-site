package postgres

import (
	"context"

	"github.com/jmoiron/sqlx"
)

type Tx interface {
	WithTx(context.Context, func(context.Context) error) error
}

type Database interface {
	sqlx.ExecerContext
	sqlx.QueryerContext
	GetContext(context.Context, interface{}, string, ...interface{}) error
	SelectContext(context.Context, interface{}, string, ...interface{}) error
}

type TxDatabase interface {
	Tx
	Database
	Do(ctx context.Context) Database
}
