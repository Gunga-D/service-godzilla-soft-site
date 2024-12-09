package service

import (
	"context"
	"fmt"
)

func Listen(ctx context.Context, listener Listener, mux Mux, opts ...Option) error {
	o := &options{
		port: _defaultPort,
	}
	for _, opt := range opts {
		opt(o)
	}
	return listen(ctx, listener, mux, o)
}

func listen(ctx context.Context, listener Listener, mux Mux, o *options) error {
	mux.Handle(_healthEndpoint, healthHandler())

	port := envInt(_envPort, o.port)
	if o.logger != nil {
		o.logger.Info(ctx, fmt.Sprintf("server started on :%d", port))
	}
	return listener.Listen(ctx, port, mux)
}
