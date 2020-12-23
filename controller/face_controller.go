package controller

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// Get faces
func GetFacesHandler(c *gin.Context) {
	faces, err := MyStore.GetFaces()

	if err == nil {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Content-Type", "application/json")
		c.JSON(http.StatusOK, faces)
	} else {
		panic(err)
	}
}

// Get a specific face
func GetFaceHandler(c *gin.Context) {
	if faceid, err := strconv.Atoi(c.Param("faceID")); err == nil {
		face, err := MyStore.GetFace(faceid)
		if err == nil {
			c.Header("Access-Control-Allow-Origin", "*")
			c.Header("Content-Type", "application/json")
			c.JSON(http.StatusOK, face)
		} else {
			panic(err)
		}
	} else {
		panic(err)
	}
}

// Run the face classifier
func ClassifyFacesHandler(c *gin.Context) {
	MyStore.ClassifyFaces()
	c.String(http.StatusOK, "Running!")
}
