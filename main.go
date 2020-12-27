package main

import (
	"fmt"
	"github.com/jschaefer-io/IDaaS/models"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"net/http"
	"os"
)

//type User struct {
//	gorm.Model
//	name     string
//	password string
//}
//
//func dbMigrate(db *gorm.DB) error {
//	return db.AutoMigrate(
//		User{},
//		User{},
//	)
//}

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

	err = orm.AutoMigrate(&models.User{})
	if err != nil {
		panic(err)
	}

	// Start application server
	srv := NewServer(orm)
	port := os.Getenv("APP_PORT")
	fmt.Printf("Starting Server on Port %s\n", port)
	err = http.ListenAndServe(":"+port, &srv)

	// panic if the http service is unable to boot
	if err != nil {
		panic(err)
	}
}
