package webserver

import (
	"net/http"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
)

type WebServer struct {
	Router        chi.Router
	Handlers      map[string]http.HandlerFunc
	WebServerPort string
}

func NewWebServer(serverPort string) *WebServer {
	return &WebServer{
		Router:        chi.NewRouter(),
		Handlers:      make(map[string]http.HandlerFunc),
		WebServerPort: serverPort,
	}
}

func (w *WebServer) Start() error {
	w.Router.Use(middleware.Logger())
	return http.ListenAndServe(w.WebServerPort, w.Router)
}

func (w *WebServer) RegisterHandler(method string, path string, handler http.HandlerFunc) {
	w.Router.MethodFunc(method, path, handler)
	w.Handlers[path] = handler
}
