package controller

import (
	"github.com/anish-krishnan/Tidepod/entity"
	"github.com/anish-krishnan/Tidepod/store"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var MyStore store.Store

func init() {
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	db.AutoMigrate(&entity.Joke{})
	db.AutoMigrate(&entity.Photo{})

	MyStore = &store.DBStore{DB: db}
}
