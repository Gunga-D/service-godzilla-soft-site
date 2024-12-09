package listener

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"
)

var (
	_defaultReadTimeout     = 10 * time.Second
	_defaultWriteTimeout    = 10 * time.Second
	_defaultIdleTimeout     = 60 * time.Second
	_defaultShutdownTimeout = 5 * time.Second
)

type HTTPListener struct {
	options
}

func NewHTTP(opts ...Option) *HTTPListener {
	o := options{
		shutdownTimeout: _defaultShutdownTimeout,
	}
	for _, opt := range opts {
		opt(&o)
	}
	return &HTTPListener{
		options: o,
	}
}

func (l *HTTPListener) Listen(ctx context.Context, port int, handler http.Handler) error {
	for _, mw := range l.mws {
		handler = mw(handler)
	}

	server := &http.Server{
		Addr:         ":" + strconv.Itoa(port),
		Handler:      handler,
		ReadTimeout:  _defaultReadTimeout,
		WriteTimeout: _defaultWriteTimeout,
		IdleTimeout:  _defaultIdleTimeout,
	}
	for _, opt := range l.onServers {
		opt(server)
	}

	chErrors := make(chan error)
	chSignals := make(chan os.Signal, 2)
	signal.Notify(chSignals, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go listen(chErrors, server)
	var err error
	select {
	case err = <-chErrors:
		shutdown(server, l.shutdownTimeout)
	case <-chSignals:
		err = shutdown(server, l.shutdownTimeout)
	case <-ctx.Done():
		err = ctx.Err()
		if e := shutdown(server, l.shutdownTimeout); e != nil {
			err = e
		}
	}
	return err
}

func listen(ch chan error, server *http.Server) {
	err := server.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		ch <- err
	}
}

func shutdown(server *http.Server, timeout time.Duration) error {
	var cancel func()
	ctx := context.Background()
	if timeout > 0 {
		ctx, cancel = context.WithTimeout(ctx, timeout)
		defer cancel()
	}

	server.SetKeepAlivesEnabled(false)
	return server.Shutdown(ctx)
}
