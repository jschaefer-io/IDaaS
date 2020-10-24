package main

import (
	"fmt"
	"github.com/jschaefer-io/IDaaS/database"

	"gorm.io/gorm"
	"net/http"
	"os"
)

type User struct {
	gorm.Model
	name     string
	password string
}

func dbMigrate(db *gorm.DB) error {
	return db.AutoMigrate(
		User{},
		User{},
	)
}

func main() {

	// Establish db connection
	db, err := database.GetConnection("mysql")
	if err != nil {
		panic(err)
	}

	// Auto-Migrate db
	err = dbMigrate(db)
	if err != nil {
		panic(err)
	}

	// Start application server
	srv := NewServer(db)
	port := os.Getenv("APP_PORT")
	fmt.Printf("Starting Server on Port %s\n", port)
	err = http.ListenAndServe(":"+port, &srv)

	// panic if the http service is unable to boot
	if err != nil {
		panic(err)
	}
}
