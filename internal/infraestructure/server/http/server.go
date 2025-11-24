package http

import (
	"context"
	"net/http"
	"time"

	"github.com/brunosprado/api-order-processor/pkg/log"
)

type Server struct {
	server *http.Server
	log    log.Logger
}

func New(port string, handler http.Handler, log log.Logger) *Server {
	return &Server{
		server: &http.Server{
			Addr:         ":" + port,
			Handler:      handler,
			ReadTimeout:  5 * time.Second,
			WriteTimeout: 55 * time.Second,
		},
		log: log,
	}
}

func (s *Server) ListenAndServe() {
	go func() {
		s.log.Info().Sendf("API Order Processor running on %s!", s.server.Addr)
		if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			s.log.Info().Sendf("Error on ListenAndServe: %q", err)
		}
	}()
}

func (s *Server) Shutdown() {
	s.log.Info().Sendf("Shutting down server")
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	if err := s.server.Shutdown(ctx); err != nil && err != http.ErrServerClosed {
		s.log.Info().Sendf("Could not shutdown in 60s: %q", err)
		return
	}
	s.log.Info().Sendf("Server gracefully stopped")
}
