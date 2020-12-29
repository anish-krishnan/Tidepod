package controller

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// GetFacesHandler gets all faces
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

// GetFaceHandler gets a specific face by ID
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

// ClassifyFacesHandler runs the face classifier on all unassigned images
func ClassifyFacesHandler(c *gin.Context) {
	MyStore.ClassifyFaces()
	c.String(http.StatusOK, "Running!")
}
