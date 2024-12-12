package amqp

import (
	"context"
	"fmt"
	"log"
	"sync"

	sdk "github.com/rabbitmq/amqp091-go"
)

var (
	once = &sync.Once{}
	amqp *sdk.Channel
)

func Get(ctx context.Context, queues []string) Amqp {
	once.Do(func() {
		host, err := loadHost()
		if err != nil {
			log.Fatalln("failed to load amqp info", err)
		}
		port, err := loadPort()
		if err != nil {
			log.Fatalln("failed to load amqp info", err)
		}
		user, err := loadUser()
		if err != nil {
			log.Fatalln("failed to load amqp info", err)
		}
		pwd, err := loadPwd()
		if err != nil {
			log.Fatalln("failed to load amqp info", err)
		}
		conn, err := sdk.Dial(fmt.Sprintf("amqp://%s:%s@%s:%d/", user, pwd, host, port))
		if err != nil {
			log.Fatalln("failed dial the connection of amqp", err)
		}
		ch, err := conn.Channel()
		if err != nil {
			log.Fatalln("failed to open amqp channel", err)
		}

		for _, queue := range queues {
			_, err = ch.QueueDeclare(
				queue, // name
				false, // durable
				false, // delete when unused
				false, // exclusive
				false, // no-wait
				nil,   // arguments
			)
			if err != nil {
				log.Fatalln("failed to declare default queue", err)
			}
		}

		go func() {
			<-ctx.Done()
			conn.Close()
			ch.Close()
		}()

		amqp = ch
	})
	return amqp
}
