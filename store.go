package main

import (
	"database/sql"
)

type Store interface {
	CreateJoke(joke *Joke) error
	GetJokes() ([]*Joke, error)
}

type dbStore struct {
	db *sql.DB
}

func (store *dbStore) CreateJoke(joke *Joke) error {
	_, err := store.db.Query("INSERT INTO jokes(likes, joke) VALUES ($1,$2)", joke.Likes, joke.Joke)
	return err
}

func (store *dbStore) GetJokes() ([]*Joke, error) {
	rows, err := store.db.Query("SELECT id, likes, joke from jokes")
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
