package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// SearchHandler handles answering a query
func SearchHandler(c *gin.Context) {
	query := c.Param("query")

	result, err := MyStore.Search(query)
	if err != nil {
		panic(err)
	}

	if err == nil {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Content-Type", "application/json")
		c.JSON(http.StatusOK, result)
	} else {
		panic(err)
	}
}
