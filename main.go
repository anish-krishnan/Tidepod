package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/contrib/static"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

/////////////////////////////////////
type Store interface {
	CreateJoke(jokeString string) error
	DeleteJoke(jokeID int) error
	LikeJoke(jokeID int) error
	GetJokes() ([]*Joke, error)
}

type dbStore struct {
	db *sql.DB
}

func (store *dbStore) CreateJoke(jokeString string) error {
	_, err := store.db.Query("INSERT INTO jokes(id, likes, joke) VALUES (nextval('jokes_id_seq'), $1, $2)", 0, jokeString)
	if err != nil {
		panic(err)
	}
	return err
}

func (store *dbStore) DeleteJoke(jokeID int) error {
	fmt.Println("ANISH DELETE", jokeID)
	_, err := store.db.Query("DELETE FROM jokes WHERE id = $1", jokeID)
	if err != nil {
		panic(err)
	}
	return err
}

func (store *dbStore) LikeJoke(jokeID int) error {
	_, err := store.db.Query("UPDATE jokes SET likes = likes + 1 WHERE id = $1", jokeID)
	if err != nil {
		panic(err)
	}
	return err
}

func (store *dbStore) GetJokes() ([]*Joke, error) {
	rows, err := store.db.Query("SELECT id, likes, joke from jokes ORDER BY id")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	jokes := []*Joke{}
	for rows.Next() {
		joke := &Joke{}

		if err := rows.Scan(&joke.ID, &joke.Likes, &joke.Joke); err != nil {
			return nil, err
		}

		jokes = append(jokes, joke)
	}
	return jokes, nil
}

var store Store

func InitStore(s Store) {
	store = s
}

/////////////////////////////////////

// Joke contains information about a single Joke
type Joke struct {
	ID    int    `json:"id" binding:"required"`
	Likes int    `json:"likes"`
	Joke  string `json:"joke" binding:"required"`
}

// We'll create a list of jokes
var jokes = []Joke{
	Joke{1, 0, "Did you hear about the restaurant on the moon? Great food, no atmosphere."},
	Joke{2, 0, "What do you call a fake noodle? An Impasta."},
	Joke{3, 0, "How many apples grow on a tree? All of them."},
	Joke{4, 0, "Want to hear a joke about paper? Nevermind it's tearable."},
	Joke{5, 0, "I just watched a program about beavers. It was the best dam program I've ever seen."},
	Joke{6, 0, "Why did the coffee file a police report? It got mugged."},
	Joke{7, 0, "How does a penguin build it's house? Igloos it together."},
}

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
	InitStore(&dbStore{db: db})

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
	// Our API will consit of just two routes
	// /jokes - which will retrieve a list of jokes a user can see
	// /jokes/like/:jokeID - which will capture likes sent to a particular joke
	api.GET("/jokes", JokeHandler)
	api.POST("/jokes/create/:joke", CreateJokeHandler)
	api.POST("/jokes/like/:jokeID", LikeJoke)
	api.POST("/jokes/delete/:jokeID", DeleteJokeHandler)

	// Start and run the server
	router.Run(":3000")
}

// JokeHandler retrieves a list of available jokes
func JokeHandler(c *gin.Context) {
	newJokes, err := store.GetJokes()

	if err == nil {
		c.Header("Content-Type", "application/json")
		c.JSON(http.StatusOK, newJokes)
	} else {
		panic(err)
	}

}

// CreateJokeHandler
func CreateJokeHandler(c *gin.Context) {
	err := store.CreateJoke(c.Param("joke"))
	if err == nil {
		c.JSON(http.StatusOK, []*Joke{})
	} else {
		panic(err)
	}
}

// DeleteJokeHandler
func DeleteJokeHandler(c *gin.Context) {
	fmt.Println("ANISH DELTE HANDLER")
	if jokeid, err := strconv.Atoi(c.Param("jokeID")); err == nil {
		err := store.DeleteJoke(jokeid)
		if err == nil {
			c.JSON(http.StatusOK, []*Joke{})
		} else {
			panic(err)
		}
	}
}

// LikeJoke increments the likes of a particular joke Item
func LikeJoke(c *gin.Context) {
	if jokeid, err := strconv.Atoi(c.Param("jokeID")); err == nil {
		err := store.LikeJoke(jokeid)
		if err == nil {
			c.JSON(http.StatusOK, []*Joke{})
		} else {
			panic(err)
		}
	}
	// if jokeid, err := strconv.Atoi(c.Param("jokeID")); err == nil {
	// 	for i := 0; i < len(jokes); i++ {
	// 		if jokes[i].ID == jokeid {
	// 			jokes[i].Likes++
	// 		}
	// 	}
	// 	c.JSON(http.StatusOK, &jokes)
	// } else {
	// 	c.AbortWithStatus(http.StatusNotFound)
	// }
}
