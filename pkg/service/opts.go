package service

type Option func(*options)

type options struct {
	port   int
	logger Logger
}

func WithPort(port int) Option {
	return func(o *options) {
		o.port = port
	}
}

func WithLogger(logger Logger) Option {
	return func(o *options) {
		o.logger = logger
	}
}
