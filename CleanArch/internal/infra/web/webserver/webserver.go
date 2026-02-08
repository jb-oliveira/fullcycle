package webserver

import (
	"net/http"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
)

type WebServer struct {
	Router        chi.Router
	WebServerPort string
}

func NewWebServer(serverPort string) *WebServer {
	router := chi.NewRouter()
	router.Use(middleware.Logger)
	return &WebServer{
		Router:        router,
		WebServerPort: serverPort,
	}
}

func (w *WebServer) Start() error {
	return http.ListenAndServe(":"+w.WebServerPort, w.Router)
}

func (w *WebServer) RegisterHandler(method string, path string, handler http.HandlerFunc) {
	w.Router.MethodFunc(method, path, handler)
}
