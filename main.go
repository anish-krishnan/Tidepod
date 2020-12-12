package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"

	"github.com/anish-krishnan/Tidepod/store"
	"github.com/anish-krishnan/Tidepod/types"

	"github.com/gin-gonic/contrib/static"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

var myStore store.Store

func main() {
	connString := "dbname=jokes_db sslmode=disable"
	db, err := sql.Open("postgres", connString)
	if err != nil {
		panic(err)
	}

	err = db.Ping()
	if err != nil {
		panic(err)
	}

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

	// Start and run the server
	router.Run(":3000")
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
		c.JSON(http.StatusOK, []*types.Joke{})
	} else {
		panic(err)
	}
}

// DeleteJokeHandler
func DeleteJokeHandler(c *gin.Context) {
	fmt.Println("ANISH DELTE HANDLER")
	if jokeid, err := strconv.Atoi(c.Param("jokeID")); err == nil {
		err := myStore.DeleteJoke(jokeid)
		if err == nil {
			c.JSON(http.StatusOK, []*types.Joke{})
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
			c.JSON(http.StatusOK, []*types.Joke{})
		} else {
			panic(err)
		}
	}
}
