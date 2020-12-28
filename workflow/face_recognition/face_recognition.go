package face_recognition

import (
	"github.com/Kagami/go-face"
	"github.com/anish-krishnan/Tidepod/entity"
	"gorm.io/gorm"
)

// ClassifyFacesByBoxEngine trains on already labelled faces, and classifies all other photos
func ClassifyFacesByBoxEngine(db *gorm.DB, boxes []*entity.Box) map[int]string {
	// Mapping photoIDs to respective boxes that are in train or test
	var trainSet []int
	var testSet []int
	var faceMap map[int]string = make(map[int]string)

	for j, box := range boxes {
		if box.Face.ID != 0 {
			trainSet = append(trainSet, j)
			faceMap[box.Face.ID] = box.Face.Name
		} else {
			testSet = append(testSet, j)
		}
	}

	rec, err := face.NewRecognizer("./workflow")
	if err != nil {
		panic(err)
	}
	defer rec.Close()

	var samples []face.Descriptor
	var labels []int32

	for _, boxIndex := range trainSet {
		face, err := rec.RecognizeSingleFile("./photo_storage/boxes/" + boxes[boxIndex].FilePath)
		if face == nil {
			continue
		} else if err != nil {
			panic(err)
		}

		samples = append(samples, face.Descriptor)
		labels = append(labels, int32(boxes[boxIndex].Face.ID))
	}

	rec.SetSamples(samples, labels)

	var result map[int]string = make(map[int]string)

	for _, boxIndex := range testSet {
		face, err := rec.RecognizeSingleFile("./photo_storage/boxes/" + boxes[boxIndex].FilePath)
		if face == nil {
			continue
		} else if err != nil {
			panic(err)
		}
		label := rec.ClassifyThreshold(face.Descriptor, 0.2)
		if label > 0 {
			result[boxes[boxIndex].ID] = faceMap[label]
		}
	}
	return result
}

// ClassifyFacesByPhotoEngine trains on already labelled faces, and classifies all other photos
func ClassifyFacesByPhotoEngine(db *gorm.DB, photos []*entity.Photo) {
	// Mapping photoIDs to respective boxes that are in train or test
	var trainSet map[int][]int = make(map[int][]int)
	var testSet map[int][]int = make(map[int][]int)

	var photoMap map[int]*entity.Photo = make(map[int]*entity.Photo)
	var faceMap map[int]string = make(map[int]string)

	for _, photo := range photos {
		for j, box := range photo.Boxes {
			photoMap[photo.ID] = photo
			if box.Face.ID != 0 {
				trainSet[photo.ID] = append(trainSet[photo.ID], j)
			} else {
				testSet[photo.ID] = append(testSet[photo.ID], j)
			}
		}
	}

	rec, err := face.NewRecognizer(".")
	if err != nil {
		panic(err)
	}
	defer rec.Close()

	var samples []face.Descriptor
	var labels []int32

	for photoID, boxes := range trainSet {

		faces, err := rec.RecognizeFile("./photo_storage/saved/" + photoMap[photoID].FilePath)
		if err != nil {
			panic(err)
		}

		for _, boxIndex := range boxes {
			samples = append(samples, faces[boxIndex].Descriptor)
			labels = append(labels, int32(photoMap[photoID].Boxes[boxIndex].Face.ID))
			faceMap[photoMap[photoID].Boxes[boxIndex].Face.ID] = photoMap[photoID].Boxes[boxIndex].Face.Name
		}
	}

	rec.SetSamples(samples, labels)

	var result map[int]string = make(map[int]string)

	for photoID, boxes := range testSet {

		for _, boxIndex := range boxes {
			face, err := rec.RecognizeSingleFile("./photo_storage/boxes/" + photoMap[photoID].Boxes[boxIndex].FilePath)
			if err != nil {
				panic(err)
			}
			if face == nil {
				continue
			}
			label := rec.ClassifyThreshold(face.Descriptor, 0.2)
			result[photoMap[photoID].Boxes[boxIndex].ID] = faceMap[label]
		}
	}
}
