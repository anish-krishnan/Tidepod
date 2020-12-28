package store

import (
	"mime/multipart"

	"github.com/anish-krishnan/Tidepod/entity"
	"gorm.io/gorm"
)

type Store interface {
	CreateJoke(jokeString string) error
	DeleteJoke(jokeID int) error
	LikeJoke(jokeID int) error
	GetJokes() ([]*entity.Joke, error)

	CreatePhoto(filename string, uploadedFile *multipart.FileHeader) error
	CreatePhotoFromMobile(filename string, uploadedFile *multipart.FileHeader, info map[string]interface{}) error
	GetPhoto(photoID int) (entity.Photo, error)
	DeletePhoto(photoID int) error
	GetPhotos() ([]*entity.Photo, error)

	GetLabels() ([]*entity.Label, error)
	GetLabel(labelID int) (entity.Label, error)

	GetFaces() ([]*entity.Face, error)
	GetFace(faceID int) (entity.Face, error)
	ClassifyFaces() error

	DeleteBox(box entity.Box) error
	AssignFaceToBox(boxID int, faceName string) (entity.Box, error)
}

type DBStore struct {
	DB *gorm.DB
}
