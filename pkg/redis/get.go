package redis

import (
	"context"
	"fmt"
	"log"
	"sync"

	redigo "github.com/gomodule/redigo/redis"
)

var (
	once = &sync.Once{}
	pool *redigo.Pool
)

func Get(ctx context.Context) *Pool {
	once.Do(func() {
		host, err := loadHost()
		if err != nil {
			log.Fatalln("failed to load db info", err)
		}
		port, err := loadPort()
		if err != nil {
			log.Fatalln("failed to load db info", err)
		}
		pwd, err := loadPwd()
		if err != nil {
			log.Fatalln("failed to load db info", err)
		}

		pool = &redigo.Pool{
			Dial: func() (redigo.Conn, error) {
				c, err := redigo.Dial("tcp", fmt.Sprintf("%s:%d", host, port))
				if err != nil {
					return nil, err
				}
				if _, err := c.Do("AUTH", pwd); err != nil {
					c.Close()
					return nil, err
				}
				return c, nil
			},
		}

		go func() {
			<-ctx.Done()
			pool.Close()
		}()
	})
	return New(pool)
}
