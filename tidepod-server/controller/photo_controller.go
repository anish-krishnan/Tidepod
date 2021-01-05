package controller

import (
	"encoding/json"
	"fmt"
	"net/http"
	"path/filepath"
	"strconv"

	"github.com/anish-krishnan/Tidepod/tidepod-server/entity"
	"github.com/anish-krishnan/Tidepod/tidepod-server/service"
	"github.com/gin-gonic/gin"
)

// GetPhotosHandler gets all photos
func GetPhotosHandler(c *gin.Context) {
	if !service.VerifyAPIRequest(c, c.Request.Header["Token"]) {
		return
	}

	photos, err := MyStore.GetPhotos()

	if err == nil {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Content-Type", "application/json")
		c.JSON(http.StatusOK, photos)
	} else {
		panic(err)
	}
}

// GetPhotosByMonthHandler gets all photos
func GetPhotosByMonthHandler(c *gin.Context) {
	if !service.VerifyAPIRequest(c, c.Request.Header["Token"]) {
		return
	}

	result, err := MyStore.GetPhotosByMonth()

	for _, pair := range result {
		fmt.Println(pair.Month, pair.Photos)
	}

	if err == nil {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Content-Type", "application/json")
		c.JSON(http.StatusOK, result)
	} else {
		panic(err)
	}
}

// GetPhotoHandler gets a specific photo by ID
func GetPhotoHandler(c *gin.Context) {
	if !service.VerifyAPIRequest(c, c.Request.Header["Token"]) {
		return
	}
	if photoid, err := strconv.Atoi(c.Param("photoID")); err == nil {
		photo, err := MyStore.GetPhoto(photoid)
		if err == nil {
			c.Header("Access-Control-Allow-Origin", "*")
			c.Header("Content-Type", "application/json")
			c.JSON(http.StatusOK, photo)
		} else {
			panic(err)
		}
	} else {
		panic(err)
	}
}

// UploadHandler handles the upload of multiple files to "photo_storage/"
// folder
func UploadHandler(c *gin.Context) {

	err := c.Request.ParseMultipartForm(1000)
	form, err := c.MultipartForm()
	if err != nil {
		c.String(http.StatusBadRequest, fmt.Sprintf("get form err: %s", err.Error()))
		fmt.Println("error getting multipartform", err)
		panic(err)
	}

	if !service.VerifyAPIRequest(c, form.Value["token"]) {
		return
	}

	files := form.File["files"]
	infoArray := form.Value["infoArray"]

	for i, file := range files {
		filename := filepath.Base(file.Filename)
		unixTime, err := strconv.ParseInt(infoArray[i], 10, 64)
		if err != nil {
			panic(err)
		}

		err = MyStore.CreatePhoto(filename, file, unixTime)
		if err != nil {
			panic(err)
		}
	}

	c.JSON(http.StatusOK, nil)
}

// PreUploadMobileHandler handles takes data about photos that the mobile
// device is trying to upload and returns the JSON containing the new files
// that should be uploaded
func PreUploadMobileHandler(c *gin.Context) {

	err := c.Request.ParseMultipartForm(1000)

	form, err := c.MultipartForm()
	if err != nil {
		c.String(http.StatusBadRequest, fmt.Sprintf("get form err: %s", err.Error()))
		fmt.Println("error getting multipartform", err)
		panic(err)
	}

	infoArray := form.Value["infoArray"]

	var result []int
	for i, info := range infoArray {

		var raw map[string]interface{}
		if err := json.Unmarshal([]byte(info), &raw); err != nil {
			panic(err)
		}

		if !MyStore.IsDuplicatePhoto(raw) {
			result = append(result, i)
		}
	}

	c.JSON(http.StatusOK, result)
}

// UploadMobileHandler handles the upload of multiple files to "photo_storage/"
// folder from a mobile device. The EXIF data arrives separately in the request
func UploadMobileHandler(c *gin.Context) {

	err := c.Request.ParseMultipartForm(1000)

	form, err := c.MultipartForm()
	if err != nil {
		c.String(http.StatusBadRequest, fmt.Sprintf("get form err: %s", err.Error()))
		fmt.Println("error getting multipartform", err)
		panic(err)
	}

	files := form.File["files"]
	infoArray := form.Value["infoArray"]

	for i, file := range files {
		filename := filepath.Base(file.Filename)
		info := infoArray[i]

		var raw map[string]interface{}
		if err := json.Unmarshal([]byte(info), &raw); err != nil {
			panic(err)
		}

		err := MyStore.CreatePhotoFromMobile(filename, file, raw)
		if err != nil {
			panic(err)
		}
	}

	// c.String(http.StatusOK, fmt.Sprintf("uploaded %d files!", len(files)))
	c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(fmt.Sprintf("uploaded %d files! <a href='http://localhost:3000'>back</a>", len(files))))
}

// DeletePhotoHandler deletes a specific photo by ID
func DeletePhotoHandler(c *gin.Context) {
	if !service.VerifyAPIRequest(c, c.Request.Header["Token"]) {
		return
	}
	if photoid, err := strconv.Atoi(c.Param("photoID")); err == nil {
		err := MyStore.DeletePhoto(photoid)
		if err == nil {
			c.JSON(http.StatusOK, []*entity.Joke{})
		} else {
			panic(err)
		}
	}
}
