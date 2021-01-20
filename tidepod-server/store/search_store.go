package store

import (
	"strings"

	"github.com/anish-krishnan/Tidepod/tidepod-server/entity"
	"gorm.io/gorm/clause"
)

// SearchResult represents the pair (month, photos taken on month)
type SearchResult struct {
	Labels []*entity.Label
	Faces  []*entity.Face
	Photos []*entity.Photo
}

func getUniqueLabels(store *DBStore) ([]*entity.Label, error) {
	var labels []*entity.Label
	store.DB.Preload(clause.Associations).Find(&labels)

	var nonEmptyLabels []*entity.Label

	for _, label := range labels {
		if len(label.Photos) > 0 {
			nonEmptyLabels = append(nonEmptyLabels, label)
		}
	}

	return nonEmptyLabels, nil
}

// Search performs a search of a query entered by the user
func (store *DBStore) Search(query string) (SearchResult, error) {

	formattedQuery := strings.Join(strings.Split(query, "&"), " ")

	var result SearchResult
	var uniquePhotos map[int]*entity.Photo = make(map[int]*entity.Photo)

	// Step 1: does the query match a label?
	labels, err := getUniqueLabels(store)
	if err != nil {
		panic(err)
	}
	for _, label := range labels {
		if strings.Contains(formattedQuery, label.LabelName) {
			fullLabel, err := store.GetLabel(label.ID)
			if err != nil {
				panic(err)
			}
			result.Labels = append(result.Labels, &entity.Label{ID: label.ID, LabelName: label.LabelName})
			for _, photo := range fullLabel.Photos {
				uniquePhotos[photo.ID] = &photo
			}
		}
	}

	// Step 2: does the query match a person's name?
	var faces []*entity.Face
	store.DB.Find(&faces)
	for _, face := range faces {
		if strings.Contains(formattedQuery, face.Name) {
			fullFace, err := store.GetFace(face.ID)
			if err != nil {
				panic(err)
			}
			result.Faces = append(result.Faces, &fullFace)

			for _, box := range fullFace.Boxes {
				uniquePhotos[box.Photo.ID] = &box.Photo
			}
		}
	}

	for _, photo := range uniquePhotos {
		result.Photos = append(result.Photos, photo)
	}

	return result, nil
}
