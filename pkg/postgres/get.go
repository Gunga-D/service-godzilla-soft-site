package postgres

import (
	"context"
	"fmt"
	"log"
	"sync"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

const (
	_maxOpenConns = 10
	_maxIdleConns = 5
)

var (
	once = &sync.Once{}
	db   *sqlx.DB
)

func Get(ctx context.Context) *Postgres {
	once.Do(func() {
		host, err := loadHost()
		if err != nil {
			log.Fatalln("failed to load db info", err)
		}
		port, err := loadPort()
		if err != nil {
			log.Fatalln("failed to load db info", err)
		}
		user, err := loadUser()
		if err != nil {
			log.Fatalln("failed to load db info", err)
		}
		pwd, err := loadPwd()
		if err != nil {
			log.Fatalln("failed to load db info", err)
		}
		dbname, err := loadDBname()
		if err != nil {
			log.Fatalln("failed to load db info", err)
		}
		psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
			"password=%s dbname=%s sslmode=disable",
			host, port, user, pwd, dbname)

		db, err = sqlx.ConnectContext(ctx, "postgres", psqlInfo)
		if err != nil {
			log.Fatalln("failed to connect db", err)
		}
		err = db.Ping()
		if err != nil {
			log.Fatalln("failed to ping db", err)
		}
		db.SetMaxOpenConns(_maxOpenConns)
		db.SetMaxIdleConns(_maxIdleConns)

		go func() {
			<-ctx.Done()
			db.Close()
		}()
	})
	return New(db)
}
