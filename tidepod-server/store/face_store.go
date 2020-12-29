package store

import (
	"github.com/anish-krishnan/Tidepod/tidepod-server/entity"
	facerecognition "github.com/anish-krishnan/Tidepod/tidepod-server/workflow/face_recognition"
	"gorm.io/gorm/clause"
)

// GetFaces gets a list of all faces
func (store *DBStore) GetFaces() ([]*entity.Face, error) {
	var faces []*entity.Face
	store.DB.Find(&faces)

	return faces, nil
}

// GetFace gets all photos for a particular face
func (store *DBStore) GetFace(faceID int) (entity.Face, error) {
	var face entity.Face
	var finalFace entity.Face
	store.DB.Preload(clause.Associations).First(&face, faceID)
	finalFace.ID = face.ID
	finalFace.Name = face.Name

	for _, box := range face.Boxes {
		var tempBox entity.Box
		store.DB.Preload(clause.Associations).First(&tempBox, box.ID)

		box.Photo = tempBox.Photo
		finalFace.Boxes = append(finalFace.Boxes, box)
	}

	return finalFace, nil
}

// ClassifyFaces runs the automatic face recognition engine
func (store *DBStore) ClassifyFaces() error {
	var boxes []*entity.Box
	store.DB.Preload(clause.Associations).Find(&boxes)

	result := facerecognition.ClassifyFacesByBoxEngine(store.DB, boxes)

	for boxID, faceName := range result {
		store.AssignFaceToBox(boxID, faceName)
	}
	return nil
}
