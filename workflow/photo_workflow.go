package workflow

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

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

// LabelPhotoWorkflow takes a photo and runs it through the tensorflow object
// detection module. Updates database entry appropriately with labels
func LabelPhotoWorkflow(db *gorm.DB, photo *entity.Photo) {
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

// CreateThumbnailWorkflow takes a photo, creates a 200x200 thumbnail
// and saves it to the photo_storage/thumbnails/ directory
func CreateThumbnailWorkflow(photo *entity.Photo) {
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

// GetReadableLocationWorkflow takes a photo, and uses the geo-coords
// to find a human readable address and updates the db entry appropriately
func GetReadableLocationWorkflow(db *gorm.DB, photo *entity.Photo) {
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
