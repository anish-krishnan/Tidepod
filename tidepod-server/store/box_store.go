package store

import (
	"errors"
	"fmt"
	"os"

	"github.com/anish-krishnan/Tidepod/tidepod-server/entity"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// DeleteBox deletes the box's photo and removes the entry from the database
func (store *DBStore) DeleteBox(box entity.Box) error {
	boxFilePath := "./photo_storage/boxes/" + box.FilePath
	err := os.Remove(boxFilePath)
	if err != nil {
		panic(err)
		return err
	}
	store.DB.Delete(&entity.Box{}, box.ID)
	return nil
}

// AssignFaceToBox takes a box and the name of a person and assigns that
// face to the box
func (store *DBStore) AssignFaceToBox(boxID int, faceName string) (entity.Box, error) {
	var box entity.Box
	store.DB.Preload(clause.Associations).First(&box, boxID)
	fmt.Println("Found Box:", box)

	var newFace entity.Face
	result := store.DB.Where("name = ?", faceName).First(&newFace)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		fmt.Println("Creating new face")
		newFace.Name = faceName
		store.DB.Create(&newFace)
	}

	box.Face = newFace
	store.DB.Save(&box)

	return box, nil
}
