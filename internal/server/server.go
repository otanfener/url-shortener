package server

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/otanfener/url-shortener/internal/service/dto"
)

type service interface {
	ShortenURL(req dto.ShortenRequest) (string, error)
	RedirectURL(code string) (string, error)
}
type logger interface {
	Info(msg string, fields map[string]interface{})
	Error(err error, fields map[string]interface{})
}
type Server struct {
	service service
	server  *http.Server
	logger  logger
}

func NewServer(service service, logger logger) *Server {
	svr := &Server{
		service: service,
		logger:  logger,
	}
	svr.server = &http.Server{
		Handler: svr.Handler(),
	}
	return svr
}

func (s *Server) Open(addr string) error {
	s.server.Addr = addr
	return s.server.ListenAndServe()
}
func (s *Server) Close() error {
	return s.server.Close()
}
func (s *Server) Handler() http.Handler {
	r := chi.NewRouter()
	r.Route("/api", func(r chi.Router) {
		h := NewHandler(s.service)
		h.Routes(r)
	})
	return r
}
