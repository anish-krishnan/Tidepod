package workflow

import (
	"encoding/json"
	"fmt"
	"image"
	"image/color"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"

	"github.com/Kagami/go-face"
	"github.com/anish-krishnan/Tidepod/tidepod-server/entity"
	"github.com/anish-krishnan/Tidepod/tidepod-server/util"
	objectDetectionScript "github.com/anish-krishnan/Tidepod/tidepod-server/workflow/object_detection/scripts"
	"github.com/codingsince1985/geo-golang"
	"github.com/codingsince1985/geo-golang/mapquest/open"
	"github.com/disintegration/imaging"
	"gorm.io/gorm"
)

// Geocoder used to convert lat/long coords of images
var geocoder geo.Geocoder
var mapquestAPIKEY string

// mutex used to prevent race condition on
// parallel database writes
var mutex sync.Mutex

// Database access point
var database *gorm.DB

// Queues for various photo workflows
var faceDetectQueue chan *entity.Photo
var labelPhotoQueue chan *entity.Photo
var completedTaskQueue chan *entity.Photo

func init() {
	mapquestApiKeyRaw, err := ioutil.ReadFile("credentials/MapQuestAPIKEY.txt")
	if err != nil {
		log.Fatal(err)
	}
	mapquestAPIKEY = string(mapquestApiKeyRaw[:len(mapquestApiKeyRaw)-1])
	geocoder = open.Geocoder(string(mapquestAPIKEY))

	// buffer up to 10000 images before blocking the upload process
	faceDetectQueue = make(chan *entity.Photo, 10000)
	labelPhotoQueue = make(chan *entity.Photo, 10000)
	completedTaskQueue = make(chan *entity.Photo)

	go RunFaceDetect()
	go LabelPhoto()
	go GarbageCollector()
}

// RunPhotoWorkflow runs the workflows sequentially
func RunPhotoWorkflow(db *gorm.DB, photo *entity.Photo) {
	database = db

	createFormattedTempImage(photo.FilePath, "./photo_storage/TEMP/"+photo.FilePath)

	CreatePhotoThumbnail(db, photo)
	UpdatePhotoWithReadableLocation(db, photo)

	faceDetectQueue <- photo
	labelPhotoQueue <- photo
}

// Creates a formatted temporary image used by all workflows
func createFormattedTempImage(filename string, tempFilePath string) {
	source, err := os.Open("./photo_storage/saved/" + filename)
	if err != nil {
		panic(err)
	}
	defer source.Close()

	destination, err := os.Create(tempFilePath)
	if err != nil {
		panic(err)
	}
	_, err = io.Copy(destination, source)
	if err != nil {
		panic(err)
	}
	destination.Close()

	util.UpdatePhotoRotation(tempFilePath)
}

// LabelPhoto takes a photo and runs it through the tensorflow object
// detection module. Updates database entry appropriately with labels
func LabelPhoto() {
	for {
		// Wait until a photo is received
		photo := <-labelPhotoQueue

		// Get Labels
		labels, err := objectDetectionScript.GetLabelsWithPythonScript("./photo_storage/TEMP/" + photo.FilePath)
		if err != nil {
			panic(err)
		}

		labelEntries := []entity.Label{}
		for _, label := range labels {
			var labelEntry entity.Label
			database.Where("label_name = ?", label).First(&labelEntry)
			labelEntries = append(labelEntries, labelEntry)
		}

		mutex.Lock()
		var newPhoto entity.Photo
		database.First(&newPhoto, photo.ID)
		newPhoto.Labels = labelEntries
		database.Save(&newPhoto)
		mutex.Unlock()

		completedTaskQueue <- photo
	}
}

// CreatePhotoThumbnail takes a photo, creates a 200x200 thumbnail
// and saves it to the photo_storage/thumbnails/ directory.
// It also rotates the thumbnail as needed
func CreatePhotoThumbnail(db *gorm.DB, photo *entity.Photo) {
	photo.ThumbnailFilePath = photo.FilePath
	db.Save(photo)

	rotation := util.GetPhotoRotation("photo_storage/saved/" + photo.FilePath)

	img, err := imaging.Open("photo_storage/saved/" + photo.FilePath)

	if err != nil {
		panic(err)
	}
	thumb := imaging.Thumbnail(img, 200, 200, imaging.CatmullRom)

	if rotation != 0 {
		thumb = imaging.Rotate(thumb, rotation, color.Gray{})
	}

	err = imaging.Save(thumb, "photo_storage/thumbnails/"+photo.FilePath)
	if err != nil {
		fmt.Println("ERROR saving thumbnail", photo.FilePath)
		panic(err)
	}
}

// RunFaceDetect takes a photo, and finds all faces in the image. It creates a
// new "Box" for each found face
func RunFaceDetect() {
	for {
		photo := <-faceDetectQueue
		fileParts := strings.Split(photo.FilePath, ".")
		ext := fileParts[len(fileParts)-1]

		img, err := imaging.Open("./photo_storage/TEMP/" + photo.FilePath)
		if err != nil {
			panic(err)
		}

		rec, err := face.NewRecognizer("./workflow")
		if err != nil {
			panic(err)
		}
		defer rec.Close()

		faces, err := rec.RecognizeFile("./photo_storage/TEMP/" + photo.FilePath)
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
			database.Create(&box)
			box.FilePath = fmt.Sprintf("%d.%s", box.ID, ext)
			database.Save(&box)

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
		database.First(&newPhoto, photo.ID)
		newPhoto.Boxes = photo.Boxes
		database.Save(&newPhoto)
		mutex.Unlock()

		completedTaskQueue <- photo
	}
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

// GarbageCollector waits for all tasks associated with a photo to be
// completed before removing the TEMP file
func GarbageCollector() {
	var activeTempImages map[int]int = make(map[int]int)
	for {
		photo := <-completedTaskQueue
		activeTempImages[photo.ID]++

		// Once 2 tasks have been completed, we can safely remove the temp image
		if activeTempImages[photo.ID] >= 2 {
			err := os.Remove("./photo_storage/TEMP/" + photo.FilePath)
			if err != nil {
				panic(err)
			}

			delete(activeTempImages, photo.ID)
		}
	}
}
