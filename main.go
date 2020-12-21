package main

import (
	"net/http"

	"github.com/anish-krishnan/Tidepod/controller"
	"github.com/gin-gonic/contrib/static"
	"github.com/gin-gonic/gin"
)

func main() {

	// Set up routes
	router := gin.Default()

	// Serve frontend static files
	router.Use(static.Serve("/", static.LocalFile("./views", true)))
	router.Static("/photo_storage", "./photo_storage")

	// Setup route group for the API
	api := router.Group("/api")
	{
		api.GET("/", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"message": "pong",
			})
		})
	}

	// Jokes Routes
	api.GET("/jokes", controller.GetJokesHandler)
	api.POST("/jokes/create/:joke", controller.CreateJokeHandler)
	api.POST("/jokes/like/:jokeID", controller.LikeJoke)
	api.POST("/jokes/delete/:jokeID", controller.DeleteJokeHandler)

	// Photos Routes
	api.GET("/photos", controller.GetPhotosHandler)
	api.POST("/photos/delete/:photoID", controller.DeletePhotoHandler)
	api.POST("/upload", controller.UploadHandler)

	// Labels Routes
	api.GET("/labels", controller.GetLabelsHandler)
	api.GET("/label/:labelID", controller.GetLabelHandler)

	// Start and run the server
	router.Run(":3000")
}
