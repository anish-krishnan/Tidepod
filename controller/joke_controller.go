package controller

import (
	"net/http"
	"strconv"

	"github.com/anish-krishnan/Tidepod/entity"
	"github.com/gin-gonic/gin"
)

// Get list of jokes
func GetJokesHandler(c *gin.Context) {
	jokes, err := MyStore.GetJokes()

	if err == nil {
		c.Header("Content-Type", "application/json")
		c.JSON(http.StatusOK, jokes)
	} else {
		panic(err)
	}
}

// Create a joke
func CreateJokeHandler(c *gin.Context) {
	err := MyStore.CreateJoke(c.Param("joke"))
	if err == nil {
		c.JSON(http.StatusOK, []*entity.Joke{})
	} else {
		panic(err)
	}
}

// Delete a joke with provided ID
func DeleteJokeHandler(c *gin.Context) {
	if jokeid, err := strconv.Atoi(c.Param("jokeID")); err == nil {
		err := MyStore.DeleteJoke(jokeid)
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
		err := MyStore.LikeJoke(jokeid)
		if err == nil {
			c.JSON(http.StatusOK, []*entity.Joke{})
		} else {
			panic(err)
		}
	}
}
