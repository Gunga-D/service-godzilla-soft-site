package listener

import (
	"net/http"
	"time"
)

type Option func(*options)

type options struct {
	onServers       []func(server *http.Server)
	mws             []func(handler http.Handler) http.Handler
	shutdownTimeout time.Duration
}

func WithWriteTimeout(timeout time.Duration) Option {
	return func(o *options) {
		o.onServers = append(o.onServers, func(server *http.Server) { server.WriteTimeout = timeout })
	}
}

func WithReadTimeout(timeout time.Duration) Option {
	return func(o *options) {
		o.onServers = append(o.onServers, func(server *http.Server) { server.ReadHeaderTimeout = timeout })
	}
}

func WithMW(mw func(handler http.Handler) http.Handler) Option {
	return func(o *options) {
		o.mws = append(o.mws, mw)
	}
}

func WithShutdownTimeout(timeout time.Duration) Option {
	return func(o *options) { o.shutdownTimeout = timeout }
}
