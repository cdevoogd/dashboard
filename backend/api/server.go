package api

import (
	"encoding/json"
	"errors"
	"net"
	"net/http"
	"strconv"

	"github.com/cdevoogd/dashboard/backend/dash"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

type Database interface {
	IsNotFoundError(err error) bool
	AddApplication(app *dash.ApplicationRecord) error
	GetApplication(id string) (*dash.ApplicationRecord, error)
	GetAllApplications() ([]*dash.ApplicationRecord, error)
	UpdateApplication(app *dash.ApplicationRecord) error
	DeleteApplication(id string) error
}

type Server struct {
	config  *ServerConfig
	address string
	router  *mux.Router
	logger  *zap.SugaredLogger
	db      Database
}

func NewServer(config *ServerConfig, logger *zap.SugaredLogger, db Database) (*Server, error) {
	if config == nil {
		return nil, errors.New("config is nil")
	}
	if logger == nil {
		return nil, errors.New("logger is nil")
	}
	if db == nil {
		return nil, errors.New("database is nil")
	}

	s := &Server{
		config:  config,
		address: net.JoinHostPort("", strconv.Itoa(int(config.Port))),
		router:  mux.NewRouter(),
		logger:  logger,
		db:      db,
	}
	s.configureRoutes()

	return s, nil
}

func (s *Server) Serve() error {
	s.logger.Infof("Starting API server on %q", s.address)
	return http.ListenAndServe(s.address, s.router)
}

func (s *Server) configureRoutes() {
	const (
		DELETE = http.MethodDelete
		GET    = http.MethodGet
		POST   = http.MethodPost
		PUT    = http.MethodPut
	)

	api := s.router.PathPrefix("/api").Subrouter()
	api.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(http.StatusOK) })
	api.HandleFunc("/applications", s.handleGetAllApplications).Methods(GET)
	api.HandleFunc("/application", s.handleCreateApplication).Methods(POST)
	api.HandleFunc("/application/{id}", s.handleGetApplication).Methods(GET)
	api.HandleFunc("/application/{id}", s.handleDeleteApplication).Methods(DELETE)
	api.HandleFunc("/application/{id}", s.handleUpdateApplication).Methods(PUT)
}

// writeJSON will encode the given object to JSON and attempt to write it to the response. Any
// errors will be logged using the server's logger.
func (s *Server) writeJSON(w http.ResponseWriter, obj any) {
	err := json.NewEncoder(w).Encode(obj)
	if err != nil {
		s.logger.Errorf("Error writing JSON response: %s", err)
	}
}

// handleEmptyPathParameter returns a 500 error to the client and logs out information about the
// request. If a path parameter is empty, that should mean that the router and handler are not using
// the same name for the parameter, or that the router is still routing the request when the
// parameter is missing.
func (s *Server) handleEmptyPathParameter(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "empty path parameter", http.StatusInternalServerError)
	s.logger.Errorw(
		"Encountered an empty URL parameter",
		"method", r.Method,
		"path", r.URL.EscapedPath(),
	)
}
