package controller

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// GetLabelsHandler gets all labels
func GetLabelsHandler(c *gin.Context) {
	labels, err := MyStore.GetLabels()

	if err == nil {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Content-Type", "application/json")
		c.JSON(http.StatusOK, labels)
	} else {
		panic(err)
	}
}

// GetLabelHandler gets a specific label by id
func GetLabelHandler(c *gin.Context) {

	if labelid, err := strconv.Atoi(c.Param("labelID")); err == nil {
		label, err := MyStore.GetLabel(labelid)
		if err == nil {
			c.Header("Access-Control-Allow-Origin", "*")
			c.Header("Content-Type", "application/json")
			c.JSON(http.StatusOK, label)
		} else {
			panic(err)
		}
	} else {
		panic(err)
	}
}
