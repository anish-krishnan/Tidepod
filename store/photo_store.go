package store

import (
	"fmt"
	"mime/multipart"
	"os"
	"strings"

	"github.com/anish-krishnan/Tidepod/entity"
	"github.com/anish-krishnan/Tidepod/util"
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

	var c *gin.Context
	newPhoto.OriginalFilename = filename
	fileParts := strings.Split(filename, ".")
	ext := fileParts[len(fileParts)-1]
	var newFilename string

	fmt.Println("working on", filename)

	if strings.ToLower(ext) == "png" {
		tempFilename := fmt.Sprintf("%d.%s", newPhoto.ID, ext)
		newFilename = fmt.Sprintf("%d.%s", newPhoto.ID, "jpg")
		newPhoto.FilePath = newFilename
		c.SaveUploadedFile(uploadedFile, "photo_storage/saved/"+tempFilename)
		util.ConvertPNGToJPG(tempFilename, newFilename)
	} else {
		newFilename = fmt.Sprintf("%d.%s", newPhoto.ID, ext)
		newPhoto.FilePath = newFilename
		c.SaveUploadedFile(uploadedFile, "photo_storage/saved/"+newFilename)
	}

	util.UpdatePhotoRotation(newFilename)

	file, err := os.Open("photo_storage/saved/" + newFilename)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	util.UpdatePhotoWithEXIF(&newPhoto, file)
	store.DB.Save(newPhoto)

	// Start the photo workflow in parallel
	go workflow.RunPhotoWorkflow(store.DB, &newPhoto)

	return nil
}

// CreatePhotoFromMobile functions identically to CreatePhoto but uses EXIF
// data in the 'info' json object
func (store *DBStore) CreatePhotoFromMobile(filename string, uploadedFile *multipart.FileHeader, info map[string]interface{}) error {
	var newPhoto entity.Photo
	store.DB.Create(&newPhoto)

	var c *gin.Context
	newPhoto.OriginalFilename = filename
	fileParts := strings.Split(filename, ".")
	ext := fileParts[len(fileParts)-1]
	var newFilename string

	fmt.Println("working on", filename)

	if strings.ToLower(ext) == "png" {
		tempFilename := fmt.Sprintf("%d.%s", newPhoto.ID, ext)
		newFilename = fmt.Sprintf("%d.%s", newPhoto.ID, "jpg")
		newPhoto.FilePath = newFilename
		c.SaveUploadedFile(uploadedFile, "photo_storage/saved/"+tempFilename)
		util.ConvertPNGToJPG(tempFilename, newFilename)
	} else {
		newFilename = fmt.Sprintf("%d.%s", newPhoto.ID, ext)
		newPhoto.FilePath = newFilename
		c.SaveUploadedFile(uploadedFile, "photo_storage/saved/"+newFilename)
	}

	util.UpdatePhotoRotation(newFilename)

	file, err := os.Open("photo_storage/saved/" + newFilename)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	util.UpdateMobilePhotoWithEXIF(&newPhoto, info)

	store.DB.Save(newPhoto)

	// Start the photo workflow in parallel
	go workflow.RunPhotoWorkflow(store.DB, &newPhoto)

	return nil

}

// GetPhotos gets all photos
func (store *DBStore) GetPhotos() ([]*entity.Photo, error) {
	var photos []*entity.Photo
	store.DB.Find(&photos)
	// for _, photo := range photos {
	// 	for j, box := range photo.Boxes {
	// 		var newBox entity.Box
	// 		store.DB.Preload(clause.Associations).First(&newBox, box.ID)
	// 		photo.Boxes[j] = newBox
	// 	}
	// }
	return photos, nil
}

// GetPhoto gets a specific photo by id
func (store *DBStore) GetPhoto(photoID int) (entity.Photo, error) {
	var photo entity.Photo
	store.DB.Preload(clause.Associations).First(&photo, photoID)
	for j, box := range photo.Boxes {
		var newBox entity.Box
		store.DB.Preload(clause.Associations).First(&newBox, box.ID)
		photo.Boxes[j] = newBox
	}
	return photo, nil
}

func (store *DBStore) DeletePhoto(photoID int) error {
	// Get the Photo entry to delete from filesystem first
	var photo entity.Photo
	store.DB.Preload(clause.Associations).First(&photo, photoID)

	fmt.Println("Deleting photo: ", photo)

	// Delete the boxes
	for _, box := range photo.Boxes {
		err := store.DeleteBox(box)
		if err != nil {
			panic(err)
		}
	}

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
