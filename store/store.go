package store

import (
	"os"

	"github.com/anish-krishnan/Tidepod/entity"
	"gorm.io/gorm"
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
	store.DB.Create(&entity.Photo{FilePath: filename})
	return nil
}

func (store *DBStore) GetPhotos() ([]*entity.Photo, error) {
	var photos []*entity.Photo
	store.DB.Find(&photos)
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
	store.DB.Delete(&entity.Photo{}, photoID)
	return nil
}
