package logger

import (
	"sync"

	tele "gopkg.in/telebot.v4"
)

var instance *logger = nil
var once sync.Once

type logger struct {
	telebot *tele.Bot
}

func Get() *logger {
	once.Do(func() {
		instance = &logger{}
	})
	return instance
}

func (l *logger) SetBot(b *tele.Bot) {
	l.telebot = b
}

func (l *logger) Log(msg string) {
	go l.telebot.Send(&tele.Chat{
		ID: -1002372045234,
	}, msg)
}
