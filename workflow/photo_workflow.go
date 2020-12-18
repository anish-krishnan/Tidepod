package workflow

import (
	"os"
	"time"

	"github.com/anish-krishnan/Tidepod/entity"
	"github.com/anish-krishnan/Tidepod/object_detection"
	"github.com/disintegration/imaging"
	"github.com/rwcarlsen/goexif/exif"
	"gorm.io/gorm"
)

// LabelPhotoWorkflow takes a photo and runs it through the tensorflow object
// detection module. Updates database entry appropriately with labels
func LabelPhotoWorkflow(db *gorm.DB, photo *entity.Photo) {
	time.Sleep(5)

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
	photo.Labels = labelEntries

	db.Save(&photo)
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
