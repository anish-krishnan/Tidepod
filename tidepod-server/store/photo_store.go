package store

import (
	"fmt"
	"mime/multipart"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/anish-krishnan/Tidepod/tidepod-server/entity"
	"github.com/anish-krishnan/Tidepod/tidepod-server/util"
	"github.com/anish-krishnan/Tidepod/tidepod-server/workflow"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm/clause"
)

// MonthPhotoPair represents the pair (month, photos taken on month)
type MonthPhotoPair struct {
	Month  string
	Photos []*entity.Photo
}

// CreatePhoto takes a filname to a newly updated photo and does:
//  1. gets EXIF information
//  2. labels the image using the tensorflow object detection package
//  3. adds the entry to the database
func (store *DBStore) CreatePhoto(filename string, uploadedFile *multipart.FileHeader, unixTime int64) error {
	var newPhoto entity.Photo
	store.DB.Create(&newPhoto)

	newPhoto.OriginalFilename = filename

	fileParts := strings.Split(filename, ".")
	ext := fileParts[len(fileParts)-1]
	tempFilename := fmt.Sprintf("TEMP.%s", ext)

	var c *gin.Context
	c.SaveUploadedFile(uploadedFile, "photo_storage/saved/"+tempFilename)

	// Timestamp
	newPhoto.Timestamp = time.Unix(int64(unixTime/1000), 0)

	mediaType := util.GetMediaType(tempFilename)
	newPhoto.MediaType = mediaType

	if mediaType == "photo" {
		newFilename := fmt.Sprintf("%d.%s", newPhoto.ID, "jpg")
		newPhoto.FilePath = newFilename
		util.ConvertImageToJPG(tempFilename, newFilename)
		file, err := os.Open("photo_storage/saved/" + newFilename)
		if err != nil {
			panic(err)
		}
		util.UpdatePhotoWithEXIF(&newPhoto, file)
		file.Close()
		store.DB.Save(newPhoto)
		// Start the photo workflow in parallel
		go workflow.RunPhotoWorkflow(store.DB, &newPhoto)

	} else {
		fileParts := strings.Split(filename, ".")
		ext := fileParts[len(fileParts)-1]
		newFilename := fmt.Sprintf("%d.%s", newPhoto.ID, ext)
		newPhoto.FilePath = newFilename
		err := os.Rename("photo_storage/saved/"+tempFilename, "photo_storage/saved/"+newFilename)
		if err != nil {
			panic(err)
		}
		store.DB.Save(newPhoto)
		workflow.CreateVideoThumbnail(store.DB, &newPhoto)
	}

	return nil
}

// CreatePhotoFromMobile functions identically to CreatePhoto but uses EXIF
// data in the 'info' json object
func (store *DBStore) CreatePhotoFromMobile(filename string, uploadedFile *multipart.FileHeader, info map[string]interface{}) error {
	var newPhoto entity.Photo
	store.DB.Create(&newPhoto)

	newPhoto.OriginalFilename = filename
	newFilename := fmt.Sprintf("%d.%s", newPhoto.ID, "jpg")
	newPhoto.FilePath = newFilename

	fileParts := strings.Split(filename, ".")
	ext := fileParts[len(fileParts)-1]
	tempFilename := fmt.Sprintf("TEMP.%s", ext)

	var c *gin.Context
	c.SaveUploadedFile(uploadedFile, "photo_storage/saved/"+tempFilename)
	util.ConvertImageToJPG(tempFilename, newFilename)
	// Rotation logic has been moved to the thumbnail
	// util.UpdatePhotoRotation(newFilename)

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
	store.DB.Order("timestamp desc").Find(&photos)
	return photos, nil
}

// GetPhotosByMonth gets all photos by month
func (store *DBStore) GetPhotosByMonth() ([]*MonthPhotoPair, error) {
	var photos []*entity.Photo
	store.DB.Order("timestamp desc").Find(&photos)

	var monthToPhotos map[string][]*entity.Photo = make(map[string][]*entity.Photo)
	var keys []string

	for _, photo := range photos {
		ts := photo.Timestamp
		monthYear := fmt.Sprintf("%s %d", ts.Month(), ts.Year())
		if len(monthToPhotos[monthYear]) == 0 {
			keys = append(keys, monthYear)
		}
		monthToPhotos[monthYear] = append(monthToPhotos[monthYear], photo)
	}

	layout := "January 2006"
	sort.Slice(keys, func(i, j int) bool {
		t1, err1 := time.Parse(layout, keys[i])
		t2, err2 := time.Parse(layout, keys[j])
		if err1 != nil || err2 != nil {
			fmt.Println("Sorting ERROR", err1, err2)
			panic(err1)
		}
		return t1.After(t2)
	})

	var result []*MonthPhotoPair
	for _, monthYear := range keys {
		result = append(result, &MonthPhotoPair{Month: monthYear, Photos: monthToPhotos[monthYear]})
	}

	return result, nil
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

// DeletePhoto deletes a specific photo by ID
func (store *DBStore) DeletePhoto(photoID int) error {
	// Get the Photo entry to delete from filesystem first
	var photo entity.Photo
	store.DB.Preload(clause.Associations).First(&photo, photoID)

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

// IsDuplicatePhoto uses the metadata to determine if this photo is already stored
func (store *DBStore) IsDuplicatePhoto(info map[string]interface{}) bool {
	var newPhoto entity.Photo
	newPhoto.OriginalFilename = info["name"].(string)
	util.UpdateMobilePhotoWithEXIF(&newPhoto, info["info"].(map[string]interface{}))

	var photos []*entity.Photo
	store.DB.Where("timestamp = ? AND original_filename = ?", newPhoto.Timestamp, newPhoto.OriginalFilename).Find(&photos)

	return len(photos) > 0
}
