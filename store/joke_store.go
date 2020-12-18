package store

import "github.com/anish-krishnan/Tidepod/entity"

func (store *DBStore) CreateJoke(jokeString string) error {
	store.DB.Create(&entity.Joke{Likes: 0, Joke: jokeString})
	return nil
}

func (store *DBStore) DeleteJoke(jokeID int) error {
	store.DB.Delete(&entity.Joke{}, jokeID)
	return nil
}

func (store *DBStore) LikeJoke(jokeID int) error {
	var joke entity.Joke

	store.DB.First(&joke, jokeID)
	store.DB.Model(&joke).Update("Likes", joke.Likes+1)

	return nil
}

func (store *DBStore) GetJokes() ([]*entity.Joke, error) {
	var jokes []*entity.Joke
	store.DB.Find(&jokes)
	return jokes, nil
}
