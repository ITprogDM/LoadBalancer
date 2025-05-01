package server

import (
	"context"
	"net/http"
	"time"
)

const (
	defaultReadTimeout     = 5 * time.Second
	defaultWriteTimeout    = 5 * time.Second
	defaultShutdownTimeout = 7 * time.Second
)

type Server struct {
	Server *http.Server
}

func (s *Server) RunServer(handler http.Handler, port string) error {
	s.Server = &http.Server{
		Addr:         ":" + port,
		Handler:      handler,
		ReadTimeout:  defaultReadTimeout,
		WriteTimeout: defaultWriteTimeout,
	}

	return s.Server.ListenAndServe()
}

func (s *Server) ShutdownServer() error {
	ctx, cancel := context.WithTimeout(context.Background(), defaultShutdownTimeout)
	defer cancel()

	return s.Server.Shutdown(ctx)
}
