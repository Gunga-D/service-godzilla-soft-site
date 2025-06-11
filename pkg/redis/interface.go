package redis

import (
	"context"
	redigo "github.com/gomodule/redigo/redis"
)

type Redis interface {
	Get(ctx context.Context, key string) (interface{}, error)
	MultiGet(ctx context.Context, keys []string) ([]interface{}, error)
	Set(ctx context.Context, key string, value interface{}, ttl *int64) error
	MultiSet(ctx context.Context, keys []string, values []interface{}, ttl *int64) error
	Exist(ctx context.Context, key string) (bool, error)
	Del(ctx context.Context, key string) error
	Add(ctx context.Context, key string, values []interface{}) error
	Members(ctx context.Context, key string) ([]interface{}, error)
	IsMember(ctx context.Context, key string, val interface{}) (bool, error)
	Rem(ctx context.Context, key string, vals []interface{}) error
	Execute(ctx context.Context, builder func(conn redigo.Conn) error) error
}
