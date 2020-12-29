package controller

import (
	"encoding/json"
	"fmt"
	"net/http"
	"path/filepath"
	"strconv"

	"github.com/anish-krishnan/Tidepod/entity"
	"github.com/gin-gonic/gin"
)

// GetPhotosHandler gets all photos
func GetPhotosHandler(c *gin.Context) {

	photos, err := MyStore.GetPhotos()

	if err == nil {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Content-Type", "application/json")
		c.JSON(http.StatusOK, photos)
	} else {
		panic(err)
	}
}

// GetPhotoHandler gets a specific photo by ID
func GetPhotoHandler(c *gin.Context) {
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

	form, err := c.MultipartForm()
	if err != nil {
		c.String(http.StatusBadRequest, fmt.Sprintf("get form err: %s", err.Error()))
		fmt.Println("error getting multipartform", err)
		panic(err)
	}

	files := form.File["files"]

	fmt.Println("Files", form.File)

	for _, file := range files {
		filename := filepath.Base(file.Filename)
		// if err := c.SaveUploadedFile(file, "photo_storage/temp/"+filename); err != nil {
		// 	c.String(http.StatusBadRequest, fmt.Sprintf("upload file err: %s", err.Error()))
		// 	panic(err)
		// 	return
		// }

		err := MyStore.CreatePhoto(filename, file)
		if err != nil {
			panic(err)
		}
	}

	// c.String(http.StatusOK, fmt.Sprintf("uploaded %d files!", len(files)))
	c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(fmt.Sprintf("uploaded %d files! <a href='http://localhost:3000'>back</a>", len(files))))
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

	fmt.Println("RESPONDING WITH", result)
	c.JSON(http.StatusOK, result)
}

// UploadMobileHandler handles the upload of multiple files to "photo_storage/"
// folder from a mobile device. The EXIF data arrives seperately in the request
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

		fmt.Println("\n\n\n", info, "\n\n\n")

		var raw map[string]interface{}
		if err := json.Unmarshal([]byte(info), &raw); err != nil {
			panic(err)
		}
		fmt.Println(raw)

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
	if photoid, err := strconv.Atoi(c.Param("photoID")); err == nil {
		err := MyStore.DeletePhoto(photoid)
		if err == nil {
			c.JSON(http.StatusOK, []*entity.Joke{})
		} else {
			panic(err)
		}
	}
}
