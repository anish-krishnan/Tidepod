package store

import (
	"fmt"

	"github.com/anish-krishnan/Tidepod/tidepod-server/entity"
	"gorm.io/gorm/clause"
)

// GetLabels gets a list of all labels that contain a non-zero
// number of photos
func (store *DBStore) GetLabels() ([]*entity.Label, error) {
	var labels []*entity.Label
	store.DB.Find(&labels)

	var nonEmptyLabels []*entity.Label

	for _, label := range labels {
		var tempLabel entity.Label
		store.DB.Preload(clause.Associations).First(&tempLabel, label.ID)

		if len(tempLabel.Photos) > 0 {
			label.Photos = append(label.Photos, tempLabel.Photos[0])
			nonEmptyLabels = append(nonEmptyLabels, label)
		}
	}
	fmt.Println(nonEmptyLabels)
	return nonEmptyLabels, nil
}

// GetLabel gets a specific labelID and returns information
// about the label and all photos tagged with that label
func (store *DBStore) GetLabel(labelID int) (entity.Label, error) {
	var label entity.Label
	var finalLabel entity.Label
	store.DB.Preload(clause.Associations).First(&label, labelID)
	finalLabel.LabelName = label.LabelName

	for _, photo := range label.Photos {
		var tempPhoto entity.Photo
		store.DB.Preload(clause.Associations).First(&tempPhoto, photo.ID)

		photo.Labels = tempPhoto.Labels
		finalLabel.Photos = append(finalLabel.Photos, tempPhoto)
	}

	return finalLabel, nil
}
