package controller

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

func AssignFaceToBoxHandler(c *gin.Context) {
	valueParts := strings.Split(c.Param("boxIDandName"), "+")
	boxid, err := strconv.Atoi(valueParts[0])

	if err != nil {
		panic(err)
	}

	faceName := valueParts[1]

	box, err := MyStore.AssignFaceToBox(boxid, faceName)
	if err != nil {
		panic(err)
	}
	c.Header("Access-Control-Allow-Origin", "*")
	c.Header("Content-Type", "application/json")
	c.JSON(http.StatusOK, box)

}
