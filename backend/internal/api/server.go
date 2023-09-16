package api

import (
	"errors"
	"net"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

type Server struct {
	config  *ServerConfig
	address string
	router  *chi.Mux
	logger  *zap.SugaredLogger
}

func NewServer(config *ServerConfig, logger *zap.SugaredLogger) (*Server, error) {
	if config == nil {
		return nil, errors.New("config is nil")
	}
	if logger == nil {
		return nil, errors.New("logger is nil")
	}

	s := &Server{
		config:  config,
		address: net.JoinHostPort("", strconv.Itoa(int(config.Port))),
		router:  chi.NewRouter(),
		logger:  logger,
	}
	s.configureRoutes()

	return s, nil
}

func (s *Server) Serve() error {
	s.logger.Infof("Starting API server on %q", s.address)
	return http.ListenAndServe(s.address, s.router)
}

func (s *Server) configureRoutes() {
	s.router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hello!"))
	})
}
