package main

import (
	"fmt"
	"net/http"
	"path/filepath"
	"strconv"

	"github.com/anish-krishnan/Tidepod/entity"
	"github.com/anish-krishnan/Tidepod/store"
	"github.com/gin-gonic/contrib/static"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var myStore store.Store

func main() {
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	db.AutoMigrate(&entity.Joke{})

	myStore = &store.DBStore{DB: db}

	// Set the router as the default one shipped with Gin
	router := gin.Default()

	// Serve frontend static files
	router.Use(static.Serve("/", static.LocalFile("./views", true)))

	// Setup route group for the API
	api := router.Group("/api")
	{
		api.GET("/", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"message": "pong",
			})
		})
	}

	/*
	 *	Our API has 4 routes
	 *	/jokes - which will retrieve a list of jokes a user can see
	 *	/jokes/create/:joke - create a new joke in the database
	 *	/jokes/like/:jokeID - which will capture likes sent to a particular joke
	 *	/jokes/delete/:jokeID - deletes a joke by ID
	 */
	api.GET("/jokes", JokeHandler)
	api.POST("/jokes/create/:joke", CreateJokeHandler)
	api.POST("/jokes/like/:jokeID", LikeJoke)
	api.POST("/jokes/delete/:jokeID", DeleteJokeHandler)

	api.POST("/upload", UploadHandler)

	// Start and run the server
	router.Run(":3000")
}

// Uploads multiple files to "saved/" folder
func UploadHandler(c *gin.Context) {
	form, err := c.MultipartForm()
	if err != nil {
		c.String(http.StatusBadRequest, fmt.Sprintf("get form err: %s", err.Error()))
		panic(err)
	}

	files := form.File["files"]

	for _, file := range files {
		filename := filepath.Base(file.Filename)
		if err := c.SaveUploadedFile(file, "saved/"+filename); err != nil {
			panic(err)
			c.String(http.StatusBadRequest, fmt.Sprintf("upload file err: %s", err.Error()))
			return
		}
	}

	c.String(http.StatusOK, fmt.Sprintf("uploaded %d files!", len(files)))
}

// JokeHandler retrieves a list of available jokes
func JokeHandler(c *gin.Context) {
	newJokes, err := myStore.GetJokes()

	if err == nil {
		c.Header("Content-Type", "application/json")
		c.JSON(http.StatusOK, newJokes)
	} else {
		panic(err)
	}
}

// CreateJokeHandler
func CreateJokeHandler(c *gin.Context) {
	err := myStore.CreateJoke(c.Param("joke"))
	if err == nil {
		c.JSON(http.StatusOK, []*entity.Joke{})
	} else {
		panic(err)
	}
}

// DeleteJokeHandler
func DeleteJokeHandler(c *gin.Context) {
	if jokeid, err := strconv.Atoi(c.Param("jokeID")); err == nil {
		err := myStore.DeleteJoke(jokeid)
		if err == nil {
			c.JSON(http.StatusOK, []*entity.Joke{})
		} else {
			panic(err)
		}
	}
}

// LikeJoke increments the likes of a particular joke Item
func LikeJoke(c *gin.Context) {
	if jokeid, err := strconv.Atoi(c.Param("jokeID")); err == nil {
		err := myStore.LikeJoke(jokeid)
		if err == nil {
			c.JSON(http.StatusOK, []*entity.Joke{})
		} else {
			panic(err)
		}
	}
}
