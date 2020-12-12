package store

import (
	"database/sql"
	"fmt"

	"github.com/anish-krishnan/Tidepod/types"
)

type Store interface {
	CreateJoke(jokeString string) error
	DeleteJoke(jokeID int) error
	LikeJoke(jokeID int) error
	GetJokes() ([]*types.Joke, error)
}

type DBStore struct {
	DB *sql.DB
}

func (store *DBStore) CreateJoke(jokeString string) error {
	_, err := store.DB.Query("INSERT INTO jokes(id, likes, joke) VALUES (nextval('jokes_id_seq'), $1, $2)", 0, jokeString)
	if err != nil {
		panic(err)
	}
	return err
}

func (store *DBStore) DeleteJoke(jokeID int) error {
	fmt.Println("ANISH DELETE", jokeID)
	_, err := store.DB.Query("DELETE FROM jokes WHERE id = $1", jokeID)
	if err != nil {
		panic(err)
	}
	return err
}

func (store *DBStore) LikeJoke(jokeID int) error {
	_, err := store.DB.Query("UPDATE jokes SET likes = likes + 1 WHERE id = $1", jokeID)
	if err != nil {
		panic(err)
	}
	return err
}

func (store *DBStore) GetJokes() ([]*types.Joke, error) {
	rows, err := store.DB.Query("SELECT id, likes, joke from jokes ORDER BY id")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	jokes := []*types.Joke{}
	for rows.Next() {
		joke := &types.Joke{}

		if err := rows.Scan(&joke.ID, &joke.Likes, &joke.Joke); err != nil {
			return nil, err
		}

		jokes = append(jokes, joke)
	}
	return jokes, nil
}
