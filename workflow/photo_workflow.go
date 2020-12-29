package workflow

import (
	"encoding/json"
	"fmt"
	"image"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"sync"

	"github.com/Kagami/go-face"
	"github.com/anish-krishnan/Tidepod/entity"
	"github.com/anish-krishnan/Tidepod/workflow/object_detection"
	"github.com/codingsince1985/geo-golang"
	"github.com/codingsince1985/geo-golang/mapquest/open"
	"github.com/disintegration/imaging"
	"gorm.io/gorm"
)

var geocoder geo.Geocoder
var mapquestAPIKEY string
var mutex sync.Mutex

func init() {
	mapquestApiKeyRaw, err := ioutil.ReadFile("credentials/MapQuestAPIKEY.txt")
	if err != nil {
		log.Fatal(err)
	}
	mapquestAPIKEY = string(mapquestApiKeyRaw[:len(mapquestApiKeyRaw)-1])
	geocoder = open.Geocoder(string(mapquestAPIKEY))
}

// RunPhotoWorkflow runs the workflows sequentially
func RunPhotoWorkflow(db *gorm.DB, photo *entity.Photo) {
	CreateThumbnail(photo)
	UpdatePhotoWithReadableLocation(db, photo)
	LabelPhoto(db, photo)
	RunFaceDetect(db, photo)
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

	mutex.Lock()
	var newPhoto entity.Photo
	db.First(&newPhoto, photo.ID)
	newPhoto.Labels = labelEntries
	db.Save(newPhoto)
	mutex.Unlock()
}

// CreateThumbnail takes a photo, creates a 200x200 thumbnail
// and saves it to the photo_storage/thumbnails/ directory
func CreateThumbnail(photo *entity.Photo) {
	img, err := imaging.Open("photo_storage/saved/" + photo.FilePath)
	if err != nil {
		panic(err)
	}
	thumb := imaging.Thumbnail(img, 200, 200, imaging.CatmullRom)

	err = imaging.Save(thumb, "photo_storage/thumbnails/"+photo.FilePath)
	if err != nil {
		fmt.Println("ERROR saving thumbnail", photo.FilePath)
		panic(err)
	}
}

// RunFaceDetect takes a photo, and finds all faces in the image. It creates a
// new "Box" for each found face
func RunFaceDetect(db *gorm.DB, photo *entity.Photo) {
	fileParts := strings.Split(photo.FilePath, ".")
	ext := fileParts[len(fileParts)-1]

	img, err := imaging.Open("./photo_storage/saved/" + photo.FilePath)
	if err != nil {
		panic(err)
	}

	rec, err := face.NewRecognizer("./workflow")
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

	mutex.Lock()
	var newPhoto entity.Photo
	db.First(&newPhoto, photo.ID)
	newPhoto.Boxes = photo.Boxes
	db.Save(newPhoto)
	mutex.Unlock()
}

// UpdatePhotoWithReadableLocation takes a photo, and uses the geo-coords
// to find a human readable address and updates the db entry appropriately
func UpdatePhotoWithReadableLocation(db *gorm.DB, photo *entity.Photo) {
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
