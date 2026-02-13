package webserver

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type method string

var POST method = "POST"
var GET method = "GET"

type WebServer struct {
	Router        chi.Router
	Handlers      map[string]map[method]http.HandlerFunc
	WebServerPort string
}

func NewWebServer(serverPort string) *WebServer {
	return &WebServer{
		Router:        chi.NewRouter(),
		Handlers:      make(map[string]map[method]http.HandlerFunc),
		WebServerPort: serverPort,
	}
}

func (s *WebServer) AddHandler(path string, m method, handler http.HandlerFunc) {
	if _, ok := s.Handlers[path]; !ok {
		s.Handlers[path] = make(map[method]http.HandlerFunc)
	}
	s.Handlers[path][m] = handler
}

// loop through the handlers and add them to the router
// register middeleware logger
// start the server
func (s *WebServer) Start() {
	s.Router.Use(middleware.Logger)
	for path, actionHandler := range s.Handlers {
		for m, handler := range actionHandler {
			switch m {
			case POST:
				s.Router.Post(path, handler)
			case GET:
				s.Router.Get(path, handler)
			}
		}
	}
	http.ListenAndServe(s.WebServerPort, s.Router)
}
