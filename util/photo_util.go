package util

import (
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	"image/png"
	"os"
	"time"

	"github.com/anish-krishnan/Tidepod/entity"
	"github.com/disintegration/imaging"
	"github.com/rwcarlsen/goexif/exif"
	"golang.org/x/image/tiff"
)

// UpdatePhotoWithEXIF extracts EXIF information from image. Updates database
// entry appropriately with this information
func UpdatePhotoWithEXIF(photo *entity.Photo, file *os.File) {
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

// UpdateMobilePhotoWithEXIF extracts EXIF information from the info JSON
// object. Updates database entry appropriately with this information
func UpdateMobilePhotoWithEXIF(photo *entity.Photo, info map[string]interface{}) {
	// EXIF information
	exifInfo := info["exif"].(map[string]interface{})

	// Camera Model
	if exifInfo["{TIFF}"] != nil {
		cameraModel := exifInfo["{TIFF}"].(map[string]interface{})["Model"]
		if cameraModel != nil {
			photo.CameraModel = cameraModel.(string)
		}
	}

	// Location
	if info["location"] != nil {
		lat := info["location"].(map[string]interface{})["latitude"]
		long := info["location"].(map[string]interface{})["longitude"]
		if lat != nil && long != nil {
			photo.Latitude = lat.(float64)
			photo.Longitude = long.(float64)
		}
	}

	// Timestamp
	datetime := time.Unix(int64(info["creationTime"].(float64)/1000), 0)
	photo.Timestamp = datetime

	if exifInfo["{Exif}"] != nil {
		// Focal Length
		focalLength := exifInfo["{Exif}"].(map[string]interface{})["FocalLength"]
		if focalLength != nil {
			photo.FocalLength = focalLength.(float64)
		}

		// Aperture
		apertureValue := exifInfo["{Exif}"].(map[string]interface{})["ApertureValue"]
		if apertureValue != nil {
			photo.ApertureFStop = apertureValue.(float64)
		}
	}
}

// UpdatePhotoRotation checks the photo for rotation inconsistencies
// and rotates the image appropriately
func UpdatePhotoRotation(filename string) {
	file, err := os.Open("photo_storage/saved/" + filename)
	if err != nil {
		panic(err)
	}

	x, err := exif.Decode(file)
	var rotation float64 = 0

	if err == nil {
		orientationRaw, err := x.Get("Orientation")

		if err == nil {
			orientation := orientationRaw.String()
			if orientation == "3" {
				rotation = 180
			} else if orientation == "6" {
				rotation = 270
			} else if orientation == "8" {
				rotation = 90
			}
		}

	}

	file.Close()
	if rotation != 0 {
		image, err := imaging.Open("photo_storage/saved/" + filename)
		if err != nil {
			panic(err)
		}
		rotatedImage := imaging.Rotate(image, rotation, color.Gray{})
		imaging.Save(rotatedImage, "photo_storage/saved/"+filename)
	}
}

// ConvertPNGToJPG converts a tif image to a jpg image
func ConvertPNGToJPG(pngFilename string, jpgFilename string) {
	pngImgFile, err := os.Open("photo_storage/saved/" + pngFilename)

	if err != nil {
		panic(err)
	}

	defer pngImgFile.Close()

	// create image from PNG file
	imgSrc, err := png.Decode(pngImgFile)
	if err != nil {
		panic(err)
	}

	// create new out JPEG file
	jpgImgFile, err := os.Create("photo_storage/saved/" + jpgFilename)
	if err != nil {
		panic(err)
	}

	defer jpgImgFile.Close()

	var opt jpeg.Options
	opt.Quality = 100

	// convert newImage to JPEG encoded byte and save to jpgImgFile
	// with quality = 100
	err = jpeg.Encode(jpgImgFile, imgSrc, &opt)
	if err != nil {
		panic(err)
	}

	relativePNGFilePath := "photo_storage/saved/" + pngFilename
	err = os.Remove(relativePNGFilePath)
	if err != nil {
		panic(err)
	}
}

// ConvertImageToJPG converts a tif image to a jpg image
func ConvertImageToJPG(originalFilename string, jpgFilename string) {
	mimeType := getMimeType(originalFilename)

	originalImage, err := os.Open("photo_storage/saved/" + originalFilename)

	if err != nil {
		panic(err)
	}

	defer originalImage.Close()

	// create image from file
	var imgSrc image.Image

	if mimeType == MimeTypeJpeg {
		imgSrc, err = jpeg.Decode(originalImage)
	} else if mimeType == MimeTypePng {
		imgSrc, err = png.Decode(originalImage)
	} else if mimeType == MimeTypeTiff {
		imgSrc, err = tiff.Decode(originalImage)
	} else {
		panic(fmt.Sprintf("Can't convert get image from mime type %s", mimeType))
	}

	if err != nil {
		panic(err)
	}

	// create new out JPEG file
	jpgImgFile, err := os.Create("photo_storage/saved/" + jpgFilename)
	if err != nil {
		panic(err)
	}

	defer jpgImgFile.Close()

	var opt jpeg.Options
	opt.Quality = 100

	// convert newImage to JPEG encoded byte and save to jpgImgFile
	// with quality = 100
	err = jpeg.Encode(jpgImgFile, imgSrc, &opt)
	if err != nil {
		panic(err)
	}

	relativeOriginalFilePath := "photo_storage/saved/" + originalFilename
	err = os.Remove(relativeOriginalFilePath)
	if err != nil {
		panic(err)
	}
}
