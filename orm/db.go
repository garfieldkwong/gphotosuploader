package orm

import (
	"sync"

	"gphotosuploader/orm/models"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// The DB type
type DB struct {
	Connection *gorm.DB
}

var instance *DB
var once sync.Once

// GetInstance Retrieve the singleton instance of DB
func GetInstance() *DB {
	once.Do(func() {
		instance = &DB{}
		var err error
		instance.Connection, err = gorm.Open(sqlite.Open("/db/db.sqlite3?loc=UTC"), &gorm.Config{
			Logger: logger.Default.LogMode(logger.Info),
		})
		if err != nil {
			panic("failed to connect database")
		}
		instance.Connection.AutoMigrate(&models.File{})
		if err != nil {
			panic("failed to connect database")
		}
	})

	return instance
}
