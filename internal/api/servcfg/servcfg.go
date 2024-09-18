package servcfg

import (
	"errors"
	"net/http"
	"strconv"
)

type options struct {
	port    int
	handler http.Handler
}

type Option func(opts *options) error

func WithPort(port int) Option {
	return func(opts *options) error {
		if port <= 0 {
			return errors.New("port parameter has to be more than 0")
		}
		opts.port = port
		return nil
	}
}

func WithHandler(handler http.Handler) Option {
	return func(opts *options) error {
		opts.handler = handler
		return nil
	}
}

func NewServer(addr string, opts ...Option) (*http.Server, error) {
	var options options
	var err error
	for _, opt := range opts {
		err = opt(&options)
		if err != nil {
			return nil, err
		}
	}

	var server *http.Server
	port := strconv.Itoa(options.port)
	server.Addr = addr + ":" + port

	if options.handler != nil {
		server.Handler = options.handler
	}

	return server, nil
}
