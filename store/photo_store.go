package store

import (
	"fmt"
	"mime/multipart"
	"os"
	"strings"

	"github.com/anish-krishnan/Tidepod/entity"
	"github.com/anish-krishnan/Tidepod/workflow"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm/clause"
)

// CreatePhoto takes a filname to a newly updated photo and does:
//  1. gets EXIF information
//  2. labels the image using the tensorflow object detection package
//  3. adds the entry to the database
func (store *DBStore) CreatePhoto(filename string, uploadedFile *multipart.FileHeader) error {
	var newPhoto entity.Photo

	store.DB.Create(&newPhoto)

	fileParts := strings.Split(filename, ".")
	ext := fileParts[len(fileParts)-1]
	newFilename := fmt.Sprintf("%d.%s", newPhoto.ID, ext)
	newPhoto.FilePath = newFilename

	var c *gin.Context
	c.SaveUploadedFile(uploadedFile, "photo_storage/saved/"+newFilename)

	file, err := os.Open("photo_storage/saved/" + newFilename)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	workflow.GetEXIFWorkflow(&newPhoto, file)
	store.DB.Save(newPhoto)
	go workflow.CreateThumbnailWorkflow(&newPhoto)
	go workflow.GetReadableLocationWorkflow(store.DB, &newPhoto)

	// Label image in parallel
	go workflow.LabelPhotoWorkflow(store.DB, &newPhoto)

	return nil
}

func (store *DBStore) GetPhotos() ([]*entity.Photo, error) {
	var photos []*entity.Photo
	store.DB.Preload(clause.Associations).Find(&photos)
	return photos, nil
}

func (store *DBStore) DeletePhoto(photoID int) error {
	// Get the Photo entry to delete from filesystem first
	var photo entity.Photo
	store.DB.First(&photo, photoID)

	relativeFilePath := "./photo_storage/saved/" + photo.FilePath
	err := os.Remove(relativeFilePath)
	if err != nil {
		panic(err)
		return err
	}

	relativeThumbFilePath := "./photo_storage/thumbnails/" + photo.FilePath
	err = os.Remove(relativeThumbFilePath)
	if err != nil {
		panic(err)
		return err
	}

	// Delete from DB next
	store.DB.Model(&photo).Association("Labels").Clear()
	store.DB.Delete(&entity.Photo{}, photoID)
	return nil
}
