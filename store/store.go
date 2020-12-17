package store

import (
	"fmt"
	"os"

	"github.com/anish-krishnan/Tidepod/entity"
	"github.com/anish-krishnan/Tidepod/object_detection"
	"github.com/rwcarlsen/goexif/exif"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Store interface {
	CreateJoke(jokeString string) error
	DeleteJoke(jokeID int) error
	LikeJoke(jokeID int) error
	GetJokes() ([]*entity.Joke, error)

	CreatePhoto(filename string) error
	DeletePhoto(photoID int) error
	GetPhotos() ([]*entity.Photo, error)
}

type DBStore struct {
	DB *gorm.DB
}

func (store *DBStore) CreateJoke(jokeString string) error {
	store.DB.Create(&entity.Joke{Likes: 0, Joke: jokeString})
	return nil
}

func (store *DBStore) DeleteJoke(jokeID int) error {
	store.DB.Delete(&entity.Joke{}, jokeID)
	return nil
}

func (store *DBStore) LikeJoke(jokeID int) error {
	var joke entity.Joke

	store.DB.First(&joke, jokeID)
	store.DB.Model(&joke).Update("Likes", joke.Likes+1)

	return nil
}

func (store *DBStore) GetJokes() ([]*entity.Joke, error) {
	var jokes []*entity.Joke
	store.DB.Find(&jokes)
	return jokes, nil
}

func (store *DBStore) CreatePhoto(filename string) error {
	var newPhoto entity.Photo
	newPhoto.FilePath = filename

	file, err := os.Open("saved/" + filename)
	if err != nil {
		panic(err)
	}

	exifInfo, err := exif.Decode(file)
	if err != nil {
		panic(err)
	}

	// Camera Model
	cameraModel, err := exifInfo.Get(exif.Model)
	if err == nil {
		newPhoto.CameraModel = cameraModel.String()
	}

	// Location
	lat, long, err := exifInfo.LatLong()
	fmt.Println("loc ", lat, long, err)
	if err == nil {
		newPhoto.Latitude = lat
		newPhoto.Longitude = long
	}

	// Timestamp
	tm, err := exifInfo.DateTime()
	fmt.Println("datetime ", tm, err)
	if err == nil {
		newPhoto.Timestamp = tm
	}

	// Focal Length
	focal, err := exifInfo.Get(exif.FocalLength)
	if err == nil {
		numer, denom, err := focal.Rat2(0)
		if err == nil {
			newPhoto.FocalLength = float64(numer) / float64(denom)
		}
	}

	// Aperture
	aperture, err := exifInfo.Get(exif.FNumber)
	if err == nil {
		numer, denom, err := aperture.Rat2(0)
		if err == nil {
			newPhoto.ApertureFStop = float64(numer) / float64(denom)
		}
	}

	store.DB.Create(&newPhoto)

	// Label image in parallel
	go runLabeler(store.DB, &newPhoto)

	return file.Close()
}

func runLabeler(db *gorm.DB, photo *entity.Photo) {
	// Get Labels
	labels, err := object_detection.GetLabelsForFile(photo.FilePath)
	if err != nil {
		panic(err)
	}
	labelEntries := []entity.Label{}
	for _, label := range labels {
		var labelEntry entity.Label
		db.Where("label_name = ?", label).First(&labelEntry)
		labelEntries = append(labelEntries, labelEntry)
	}
	photo.Labels = labelEntries

	db.Save(&photo)
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

	relativeFilePath := "./saved/" + photo.FilePath
	err := os.Remove(relativeFilePath)

	if err != nil {
		panic(err)
		return err
	}

	// Delete from DB next
	store.DB.Model(&photo).Association("Labels").Clear()
	store.DB.Delete(&entity.Photo{}, photoID)
	return nil
}
