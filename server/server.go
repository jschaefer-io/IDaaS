package server

import (
	"github.com/go-chi/chi/v5"
	"log"
	"net/http"
)

type Server struct {
	*Components
	Settings *Settings
}

func (s *Server) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	s.Router.ServeHTTP(writer, request)
}

type Args struct {
	Logger       *log.Logger
	Router       chi.Router
	Settings     *Settings
	Repositories *Repositories
}

func NewServer(args Args, routeInit func(*Server)) *Server {
	components := NewComponents(args.Settings, args.Repositories)
	components.Logger = args.Logger
	components.Router = args.Router

	srv := &Server{
		Components: components,
		Settings:   args.Settings,
	}
	routeInit(srv)
	return srv
}

func (s *Server) Heartbeat() {
	//for {
	//	s.Logger.Println("Server is alive")
	//	time.Sleep(time.Second * 1)
	//}
}
