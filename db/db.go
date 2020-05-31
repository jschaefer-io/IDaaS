package db

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"sync"
)

var checker = sync.Once{}

// Static db instance
var db *gorm.DB

// build a new db instance
func newDb() (*gorm.DB, error) {
	return gorm.Open("sqlite3", "./test.db")
}

// Fetches the static db instance
// or creates it, if it does not exist yet
func Get() *gorm.DB {
	checker.Do(func() {
		newDb, err := newDb()
		if err != nil {
			panic(err)
		}
		db = newDb
	})
	return db
}
