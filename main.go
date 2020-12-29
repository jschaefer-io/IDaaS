package main

import (
	"fmt"
	"github.com/jschaefer-io/IDaaS/model"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"net/http"
	"os"
)

func prepareOrm() (*gorm.DB, error) {
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASS")
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	name := os.Getenv("DB_NAME")
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", user, password, host, port, name)
	return gorm.Open(mysql.Open(dsn), &gorm.Config{})
}

func main() {
	orm, err := prepareOrm()
	if err != nil {
		panic(err)
	}

	err = orm.AutoMigrate(
		&model.User{},
		&model.Session{},
	)
	if err != nil {
		panic(err)
	}

	// create application server
	srv := NewServer(orm)

	// start application cleanup processes
	go srv.Heartbeat()

	// start application server
	port := os.Getenv("APP_PORT")
	fmt.Printf("Starting Server on Port %s\n", port)
	err = http.ListenAndServe(":"+port, &srv)

	// panic if the http service is unable to boot
	if err != nil {
		panic(err)
	}
}
