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
	return &WebServer{
		Router:        chi.NewRouter(),
		WebServerPort: serverPort,
	}
}

// Start aplica log e escuta a porta. Registre rotas com ws.Router em main antes.
func (s *WebServer) Start() error {
	s.Router.Use(middleware.Logger)
	return http.ListenAndServe(s.WebServerPort, s.Router)
}
