package webserver

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

// WebServer replica o padrão do repositório de referência (Chi + middleware de log).
type WebServer struct {
	Router        *chi.Mux
	WebServerPort string
}

// NewWebServer cria o servidor HTTP com router Chi.
func NewWebServer(serverPort string) *WebServer {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	return &WebServer{
		Router:        r,
		WebServerPort: serverPort,
	}
}

// Start escuta a porta. Middlewares já estão aplicados em NewWebServer antes das rotas.
func (s *WebServer) Start() error {
	return http.ListenAndServe(s.WebServerPort, s.Router)
}
