package controller

import (
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/anish-krishnan/Tidepod/entity"
	"github.com/anish-krishnan/Tidepod/store"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// MyStore is a global storage reference for all controllers
var MyStore store.Store

func init() {
	db := initGORM()
	migrateDatabases(db)

	MyStore = &store.DBStore{DB: db}
}

// Read labels from text file and output a string array
func readLabels(labelsFile string) ([]string, error) {
	fileBytes, err := ioutil.ReadFile(labelsFile)
	if err != nil {
		return nil, fmt.Errorf("Unable to read labels file: %v", err)
	}

	return strings.Split(string(fileBytes), "\n"), nil
}

func initGORM() *gorm.DB {
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	return db
}

func migrateDatabases(db *gorm.DB) {
	db.AutoMigrate(&entity.Joke{})
	db.AutoMigrate(&entity.Photo{})

	db.Migrator().DropTable(&entity.Label{})
	db.AutoMigrate(&entity.Label{})

	db.AutoMigrate(&entity.Face{})
	db.AutoMigrate(&entity.Box{})

	labels, err := readLabels("workflow/object_detection/ssd_mobilenet_v1_coco_2018_01_28/saved_model/labels.txt")
	if err != nil {
		panic(err)
	}

	for _, s := range labels {
		db.Create(&entity.Label{LabelName: s})
	}
}
