package controller

import (
	"net/http"
	"strconv"

	"github.com/anish-krishnan/Tidepod/entity"
	"github.com/gin-gonic/gin"
)

// GetJokesHandler gets a list of all jokes stored in the database
func GetJokesHandler(c *gin.Context) {
	jokes, err := MyStore.GetJokes()

	if err == nil {
		c.Header("Content-Type", "application/json")
		c.JSON(http.StatusOK, jokes)
	} else {
		panic(err)
	}
}

// CreateJokeHandler create a joke
func CreateJokeHandler(c *gin.Context) {
	err := MyStore.CreateJoke(c.Param("joke"))
	if err == nil {
		c.JSON(http.StatusOK, []*entity.Joke{})
	} else {
		panic(err)
	}
}

// DeleteJokeHandler deletes specific joke with provided ID
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
