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

	// router.Use(cors.New(cors.Config{
	// 	AllowedOrigins: []string{"http://localhost:3000", "http://localhost:3001"},
	// 	AllowedMethods: []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"},
	// 	AllowedHeaders: []string{"Origin", "Content-Length", "Content-Type", "Access-Control-Allow-Origin"},
	// 	ExposedHeaders: []string{"X-Total-Count"},
	// }))

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
	api.GET("/photo/:photoID", controller.GetPhotoHandler)
	api.POST("/photos/delete/:photoID", controller.DeletePhotoHandler)
	api.POST("/upload", controller.UploadHandler)

	// Labels Routes
	api.GET("/labels", controller.GetLabelsHandler)
	api.GET("/label/:labelID", controller.GetLabelHandler)

	// Faces Routes
	api.GET("/faces", controller.GetFacesHandler)
	api.GET("/face/:faceID", controller.GetFaceHandler)
	api.GET("/classifyFaces", controller.ClassifyFacesHandler)

	// Boxes Routes
	api.GET("/boxes/assignface/:boxIDandName", controller.AssignFaceToBoxHandler)

	// Start and run the server
	router.Run(":3000")
}
