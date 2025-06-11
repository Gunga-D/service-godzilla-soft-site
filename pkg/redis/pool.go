package redis

import (
	"context"

	redigo "github.com/gomodule/redigo/redis"
	"go.octolab.org/safe"
	"go.octolab.org/unsafe"
)

type Pool struct {
	conn *redigo.Pool
}

func New(conn *redigo.Pool) *Pool {
	return &Pool{
		conn: conn,
	}
}

func (c *Pool) Get(ctx context.Context, key string) (interface{}, error) {
	conn, err := c.conn.GetContext(ctx)
	if err != nil {
		return 0, err
	}
	defer func() {
		conn.Close()
	}()

	data, err := conn.Do("GET", key)
	if err == redigo.ErrNil {
		return nil, nil
	}
	return data, err
}

func (c *Pool) MultiGet(ctx context.Context, keys []string) ([]interface{}, error) {
	conn, err := c.conn.GetContext(ctx)
	if err != nil {
		return nil, err
	}

	defer func() {
		conn.Close()
	}()

	defer safe.Close(conn, unsafe.Ignore)

	return redigo.Values(conn.Do("MGET", convertKeys(keys)...))
}

func (c *Pool) Set(ctx context.Context, key string, value interface{}, ttl *int64) error {
	conn, err := c.conn.GetContext(ctx)
	if err != nil {
		return err
	}
	defer func() {
		conn.Close()
	}()

	if ttl != nil {
		_, err = conn.Do("SETEX", key, *ttl, value)
	} else {
		_, err = conn.Do("SET", key, value)
	}

	return err
}

func (c *Pool) MultiSet(ctx context.Context, keys []string, values []interface{}, ttl *int64) error {
	conn, err := c.conn.GetContext(ctx)
	if err != nil {
		return err
	}
	defer func() {
		conn.Close()
	}()

	err = conn.Send("MULTI")
	if err != nil {
		return err
	}
	for i, key := range keys {
		value := values[i]
		_ = conn.Send("SET", key, value)
		if ttl != nil {
			_ = conn.Send("EXPIRE", key, *ttl)
		}
	}
	num, err := conn.Do("EXEC")
	if num != 0 {
		return nil
	}
	return err
}

func (c *Pool) Exist(ctx context.Context, key string) (bool, error) {
	conn, err := c.conn.GetContext(ctx)
	if err != nil {
		return false, err
	}
	defer func() {
		conn.Close()
	}()

	return redigo.Bool(conn.Do("EXISTS", key))
}

func (c *Pool) Del(ctx context.Context, key string) error {
	conn, err := c.conn.GetContext(ctx)
	if err != nil {
		return err
	}
	defer func() {
		conn.Close()
	}()

	_, err = conn.Do("DEL", key)

	return err
}

func (c *Pool) Add(ctx context.Context, key string, values []interface{}) error {
	conn, err := c.conn.GetContext(ctx)
	if err != nil {
		return err
	}
	defer func() {
		conn.Close()
	}()
	args := []interface{}{key}
	args = append(args, values...)
	_, err = conn.Do("SADD", args...)
	return err
}

func (c *Pool) Members(ctx context.Context, key string) ([]interface{}, error) {
	conn, err := c.conn.GetContext(ctx)
	if err != nil {
		return nil, err
	}
	defer func() {
		conn.Close()
	}()
	return redigo.Values(conn.Do("SMEMBERS", key))
}

func (c *Pool) IsMember(ctx context.Context, key string, val interface{}) (bool, error) {
	conn, err := c.conn.GetContext(ctx)
	if err != nil {
		return false, err
	}
	defer func() {
		conn.Close()
	}()

	return redigo.Bool(conn.Do("SISMEMBER", key, val))
}

func (c *Pool) Rem(ctx context.Context, key string, vals []interface{}) error {
	conn, err := c.conn.GetContext(ctx)
	if err != nil {
		return err
	}
	defer func() {
		conn.Close()
	}()

	args := []interface{}{key}
	args = append(args, vals...)
	_, err = conn.Do("SREM", args...)

	return err
}

func convertKeys(keys []string) []interface{} {
	args := make([]interface{}, 0, len(keys))
	for _, key := range keys {
		args = append(args, key)
	}

	return args
}

func (c *Pool) Execute(ctx context.Context, executor func(conn redigo.Conn) error) error {
	conn, err := c.conn.GetContext(ctx)
	if err != nil {
		return err
	}
	defer func() {
		conn.Close()
	}()

	return executor(conn)
}
