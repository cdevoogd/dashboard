package dashboard

import (
	"bytes"
	"errors"
	"fmt"
	"html/template"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/cdevoogd/dashboard/assets"
)

// Server is the HTTP server that listens for and handles requests to the configured port.
type Server struct {
	config            *Config
	logger            *slog.Logger
	dashboardTemplate *template.Template
}

// NewServer will instantiate a new Server struct using the given config and logger. Any nil
// dependencies will result in an error being returned.
func NewServer(config *Config, logger *slog.Logger) (*Server, error) {
	if config == nil {
		return nil, errors.New("config is nil")
	}
	if logger == nil {
		return nil, errors.New("logger is nil")
	}

	dashboardTemplate, err := template.New("index.html").
		Funcs(template.FuncMap{"stripURLScheme": stripURLScheme}).
		ParseFS(assets.TemplateFS, "templates/index.html")
	if err != nil {
		return nil, fmt.Errorf("failed to parse dashboard template: %w", err)
	}

	return &Server{
		config:            config,
		logger:            logger,
		dashboardTemplate: dashboardTemplate,
	}, nil
}

// ListenAndServe will being listening for HTTP requests on the configured port.
func (s *Server) ListenAndServe() error {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /{$}", s.handleGetDashboard)
	mux.HandleFunc("GET /health", s.handleHealthCheck)

	assetServer := http.FileServer(http.FS(assets.PublicAssetFS))
	assetHandler := http.StripPrefix("/assets/", assetServer)
	mux.Handle("GET /assets/", assetHandler)

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", s.config.Port),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
		Handler:      mux,
	}
	return srv.ListenAndServe()
}

func (s *Server) handleGetDashboard(w http.ResponseWriter, r *http.Request) {
	buf := &bytes.Buffer{}
	err := s.dashboardTemplate.Execute(buf, s.config)
	if err != nil {
		e := fmt.Errorf("error executing dashboard template: %w", err)
		s.handleError(w, http.StatusInternalServerError, e)
		return
	}

	s.handleWrite(w, buf.Bytes())
}

func (s *Server) handleHealthCheck(w http.ResponseWriter, r *http.Request) {
	s.handleWrite(w, []byte("ok"))
}

func (s *Server) handleWrite(w http.ResponseWriter, content []byte) {
	_, err := w.Write(content)
	if err != nil {
		s.logger.Error("Error writing response", "err", err)
	}
}

func (s *Server) handleError(w http.ResponseWriter, code int, err error) {
	s.logger.Error("server error", "code", code, "err", err)
	http.Error(w, err.Error(), code)
}

func stripURLScheme(url string) string {
	parts := strings.SplitN(url, "://", 2)
	switch len(parts) {
	// URL might not have had a scheme
	case 1:
		return parts[0]
	// Index 0 should be the scheme, index 1 should be the rest of the URL
	case 2:
		return parts[1]
	default:
		return url
	}
}
