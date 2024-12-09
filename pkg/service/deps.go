package service

import (
	"context"
	"net/http"
)

type Listener interface {
	Listen(context.Context, int, http.Handler) error
}

type Mux interface {
	http.Handler
	Handle(string, http.Handler)
}

type Logger interface {
	Info(ctx context.Context, args ...interface{})
}
