package store

import (
	"os"

	"github.com/anish-krishnan/Tidepod/entity"
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
