package database

import (
	"errors"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"os"
)

func GetConnection(connectionType string) (*gorm.DB, error) {

	var db *gorm.DB
	var err error

	switch connectionType {
	case "mysql":
		db, err = getMySqlConnection()
		break
	default:
		return nil, errors.New(fmt.Sprintf("unsupported connection type %s", connectionType))
	}

	if err != nil {
		return nil, errors.New(fmt.Sprintf("unable to connect to database: %s", connectionType))
	}
	return db, nil
}

func getMySqlConnection() (*gorm.DB, error) {
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASS")
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	name := os.Getenv("DB_NAME")
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", user, password, host, port, name)
	return gorm.Open(mysql.Open(dsn), &gorm.Config{})
}
