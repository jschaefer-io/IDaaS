package main

import (
	"database/sql"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/jschaefer-io/IDaaS/repository"
	"github.com/jschaefer-io/IDaaS/server"
	_ "github.com/lib/pq"
)

func main() {
	logger := log.Default()

	// Settings
	settings, err := server.SettingsFromEnv()
	if err != nil {
		panic(err)
	}

	// Establish DB connection
	db, err := sql.Open("postgres", fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		os.Getenv("JIO_DB_USER"),
		os.Getenv("JIO_DB_PASSWORD"),
		os.Getenv("JIO_DB_HOST"),
		os.Getenv("JIO_DB_PORT"),
		os.Getenv("JIO_DB_NAME"),
	))
	if err != nil {
		panic(err)
	}

	// create repositories
	repos := &server.Repositories{
		UserRepository:         repository.NewUserRepository(db),
		RefreshChainRepository: repository.NewRefreshTokenRepository(db),
	}

	// create main server instance
	srv := server.NewServer(server.Args{
		Logger:       logger,
		Router:       chi.NewRouter(),
		Settings:     settings,
		Repositories: repos,
	}, ServerRoutes)

	// start application server
	host := "0.0.0.0:8080"
	listener, err := net.Listen("tcp", host)
	if err != nil {
		panic(err)
	}
	logger.Printf("server ready to accept connections on %s", host)
	err = http.Serve(listener, srv)
	if err != nil {
		panic(err)
	}
}
