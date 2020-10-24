package main

import (
	"github.com/go-chi/chi"
	"gorm.io/gorm"
	"net/http"
)

type Server struct {
	db     *gorm.DB
	router *chi.Mux
}


func NewServer(db *gorm.DB) Server{
	srv:=Server{
		db: db,
	}
	srv.init()
	return srv
}

func (s *Server) ServeHTTP(writer http.ResponseWriter, r *http.Request){
	s.router.ServeHTTP(writer, r)
}


