package workflow

import (
	"encoding/json"
	"fmt"
	"image"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/Kagami/go-face"
	"github.com/anish-krishnan/Tidepod/entity"
	"github.com/anish-krishnan/Tidepod/object_detection"
	"github.com/codingsince1985/geo-golang"
	"github.com/codingsince1985/geo-golang/mapquest/open"
	"github.com/disintegration/imaging"
	"github.com/rwcarlsen/goexif/exif"
	"gorm.io/gorm"
)

var geocoder geo.Geocoder
var mapquestAPIKEY string

func init() {
	mapquestApiKeyRaw, err := ioutil.ReadFile("credentials/MapQuestAPIKEY.txt")
	if err != nil {
		log.Fatal(err)
	}
	mapquestAPIKEY = string(mapquestApiKeyRaw[:len(mapquestApiKeyRaw)-1])
	geocoder = open.Geocoder(string(mapquestAPIKEY))
}

// GetEXIFWorkflow extracts EXIF information from image. Updates database
// entry appropriately with this information
func GetEXIFWorkflow(photo *entity.Photo, file *os.File) {
	exifInfo, err := exif.Decode(file)
	if err != nil {
		return
	}

	// Camera Model
	cameraModel, err := exifInfo.Get(exif.Model)
	if err == nil {
		photo.CameraModel = cameraModel.String()
	}

	// Location
	lat, long, err := exifInfo.LatLong()
	if err == nil {
		photo.Latitude = lat
		photo.Longitude = long
	}

	// Timestamp
	tm, err := exifInfo.DateTime()
	if err == nil {
		photo.Timestamp = tm
	}

	// Focal Length
	focal, err := exifInfo.Get(exif.FocalLength)
	if err == nil {
		numer, denom, err := focal.Rat2(0)
		if err == nil {
			photo.FocalLength = float64(numer) / float64(denom)
		}
	}

	// Aperture
	aperture, err := exifInfo.Get(exif.FNumber)
	if err == nil {
		numer, denom, err := aperture.Rat2(0)
		if err == nil {
			photo.ApertureFStop = float64(numer) / float64(denom)
		}
	}
}

// RunPhotoWorkflow runs the workflows sequentially
func RunPhotoWorkflow(db *gorm.DB, photo *entity.Photo) {
	CreateThumbnail(photo)
	GetReadableLocation(db, photo)
	LabelPhoto(db, photo)
	GetFaces(db, photo)
}

// LabelPhoto takes a photo and runs it through the tensorflow object
// detection module. Updates database entry appropriately with labels
func LabelPhoto(db *gorm.DB, photo *entity.Photo) {
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
	// photo.Labels = labelEntries

	var newPhoto entity.Photo
	db.First(&newPhoto, photo.ID)
	newPhoto.Labels = labelEntries
	db.Save(newPhoto)
}

// CreateThumbnail takes a photo, creates a 200x200 thumbnail
// and saves it to the photo_storage/thumbnails/ directory
func CreateThumbnail(photo *entity.Photo) {
	// use all CPU cores for maximum performance
	// runtime.GOMAXPROCS(runtime.NumCPU())

	img, err := imaging.Open("photo_storage/saved/" + photo.FilePath)
	if err != nil {
		panic(err)
	}
	thumb := imaging.Thumbnail(img, 200, 200, imaging.CatmullRom)

	err = imaging.Save(thumb, "photo_storage/thumbnails/"+photo.FilePath)
	if err != nil {
		panic(err)
	}
}

// GetFaces takes a photo, and finds all faces in the image. It creates a
// new "Box" for each found face
func GetFaces(db *gorm.DB, photo *entity.Photo) {
	fileParts := strings.Split(photo.FilePath, ".")
	ext := fileParts[len(fileParts)-1]

	img, err := imaging.Open("./photo_storage/saved/" + photo.FilePath)
	if err != nil {
		panic(err)
	}

	rec, err := face.NewRecognizer(".")
	if err != nil {
		panic(err)
	}
	defer rec.Close()

	faces, err := rec.RecognizeFile("./photo_storage/saved/" + photo.FilePath)
	if err != nil {
		panic(err)
	}

	for _, face := range faces {
		box := entity.Box{
			PhotoID: photo.ID,
			MinX:    face.Rectangle.Min.X,
			MinY:    face.Rectangle.Min.Y,
			MaxX:    face.Rectangle.Max.X,
			MaxY:    face.Rectangle.Max.Y,
		}
		db.Create(&box)
		box.FilePath = fmt.Sprintf("%d.%s", box.ID, ext)
		db.Save(&box)

		photo.Boxes = append(photo.Boxes, box)
		// crop out a rectangular region
		croppedImg := imaging.Crop(img, image.Rect(box.MinX, box.MinY, box.MaxX, box.MaxY))
		// lower resolution
		compresedCroppedImg := imaging.Thumbnail(croppedImg, 100, 100, imaging.CatmullRom)
		// save cropped image
		err = imaging.Save(compresedCroppedImg, "./photo_storage/boxes/"+box.FilePath)
		if err != nil {
			panic(err)
		}
	}
	var newPhoto entity.Photo
	db.First(&newPhoto, photo.ID)
	newPhoto.Boxes = photo.Boxes
	db.Save(newPhoto)
}

// ClassifyFaces trains on already labelled faces, and classifies all other photos
func ClassifyFacesEngine(db *gorm.DB, photos []*entity.Photo) {
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

	fmt.Println("TRAIN", trainSet)
	fmt.Println("TEST", testSet)

	rec, err := face.NewRecognizer(".")
	if err != nil {
		panic(err)
	}
	defer rec.Close()

	var samples []face.Descriptor
	var labels []int32

	for photoID, boxes := range trainSet {
		fmt.Println("photoid:", photoID, "=>", "boxes:", boxes)
		fmt.Println(photoMap[photoID])

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

	fmt.Println("SAMPLES", samples)
	fmt.Println("LABELS", labels)
	rec.SetSamples(samples, labels)

	var result map[int]string = make(map[int]string)

	for photoID, boxes := range testSet {
		fmt.Println("photoid:", photoID, "=>", "boxes:", boxes)

		// faces, err := rec.RecognizeFile("./photo_storage/saved/" + photoMap[photoID].FilePath)
		// if err != nil {
		// 	panic(err)
		// }

		for _, boxIndex := range boxes {
			fmt.Println("./photo_storage/boxes/" + photoMap[photoID].Boxes[boxIndex].FilePath)
			face, err := rec.RecognizeSingleFile("./photo_storage/boxes/" + photoMap[photoID].Boxes[boxIndex].FilePath)
			if err != nil {
				panic(err)
			}
			if face == nil {
				continue
			}
			label := rec.ClassifyThreshold(face.Descriptor, 0.2)
			fmt.Println("CLASSIFED", photoID, boxIndex, label)
			result[photoMap[photoID].Boxes[boxIndex].ID] = faceMap[label]
		}
	}
	fmt.Println("RESULT", result)
}

// GetReadableLocation takes a photo, and uses the geo-coords
// to find a human readable address and updates the db entry appropriately
func GetReadableLocation(db *gorm.DB, photo *entity.Photo) {
	if photo.Latitude == 0.0 && photo.Longitude == 0.0 {
		return
	}

	request := fmt.Sprintf("https://www.mapquestapi.com/geocoding/v1/reverse?key=%s&location=%f%%2C%f&outFormat=json&thumbMaps=false",
		mapquestAPIKEY,
		photo.Latitude,
		photo.Longitude)

	resp, err := http.Get(request)
	if err != nil {
		panic(err)
	}

	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)
	location := result["results"].([]interface{})[0].(map[string]interface{})["locations"].([]interface{})[0].(map[string]interface{})

	street := location["street"]
	neighborhood := location["adminArea6"]
	city := location["adminArea5"]
	postalCode := location["postalCode"]
	county := location["adminArea4"]
	state := location["adminArea3"]
	country := location["adminArea1"]

	formattedAddress := fmt.Sprintf("%s, %s %s, %s, %s, %s, %s", street, neighborhood, city, county, postalCode, state, country)
	db.Model(&photo).Updates(entity.Photo{LocationString: formattedAddress})
}
