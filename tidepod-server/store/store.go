package store

import (
	"mime/multipart"

	"github.com/anish-krishnan/Tidepod/tidepod-server/entity"
	"gorm.io/gorm"
)

// Store is an interface defining all methods to interact with the models:
//	- Joke
//	- Photo
//	- Label
//	- Face
//	- Box
type Store interface {
	CreateJoke(jokeString string) error
	DeleteJoke(jokeID int) error
	LikeJoke(jokeID int) error
	GetJokes() ([]*entity.Joke, error)

	CreatePhoto(filename string, uploadedFile *multipart.FileHeader, unixTime int64) error
	CreatePhotoFromMobile(filename string, uploadedFile *multipart.FileHeader, info map[string]interface{}) error
	GetPhoto(photoID int) (entity.Photo, error)
	DeletePhoto(photoID int) error
	GetPhotos() ([]*entity.Photo, error)
	GetPhotosByMonth(offset int) ([]*MonthPhotoPair, error)
	IsDuplicatePhoto(info map[string]interface{}) bool

	GetLabels() ([]*entity.Label, error)
	GetLabel(labelID int) (entity.Label, error)

	GetFaces() ([]*entity.Face, error)
	GetFace(faceID int) (entity.Face, error)
	ClassifyFaces() error

	GetUnassignedBoxes() ([]*entity.Box, error)
	DeleteBox(box entity.Box) error
	AssignFaceToBox(boxID int, faceName string) (entity.Box, error)

	Search(query string) (SearchResult, error)
}

// DBStore implements the Store interface
type DBStore struct {
	DB *gorm.DB
}
