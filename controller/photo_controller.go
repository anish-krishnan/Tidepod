package controller

import (
	"fmt"
	"net/http"
	"path/filepath"
	"strconv"

	"github.com/anish-krishnan/Tidepod/entity"
	"github.com/gin-gonic/gin"
)

// Get photos
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

// Uploads multiple files to "photo_storage/" folder
func UploadHandler(c *gin.Context) {

	form, err := c.MultipartForm()
	if err != nil {
		c.String(http.StatusBadRequest, fmt.Sprintf("get form err: %s", err.Error()))
		fmt.Println("error getting multipartform", err)
		panic(err)
	}

	files := form.File["files"]

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

// Delete a photo
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
