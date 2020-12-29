package main

import (
	"github.com/go-chi/chi"
	"github.com/jschaefer-io/IDaaS/model"
	"gorm.io/gorm"
	"net/http"
	"time"
)

type Server struct {
	db     *gorm.DB
	router *chi.Mux
}

func NewServer(db *gorm.DB) Server {
	srv := Server{
		db: db,
	}
	srv.init()
	return srv
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

func (s *Server) Heartbeat() {
	for {

		// delete expired session
		s.db.Where("expires_at < ?", time.Now()).Delete(model.Session{})

		// sleep until next heartbeat tick
		time.Sleep(time.Minute)
	}
}
