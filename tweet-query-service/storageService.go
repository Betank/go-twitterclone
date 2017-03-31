package main

import "sync"

type Storage interface {
	GetTweetById(id string) tweet
	GetAllTweets() []tweet
	CreateTweet(tweet tweet)
	DeleteTweet(id string)
}

type simpleStore struct {
	sync.Mutex
	tweetStorage map[string]tweet
}

func (store *simpleStore) GetTweetById(id string) tweet {
	store.Lock()
	defer store.Unlock()
	tweet := store.tweetStorage[id]
	return tweet
}

func (store *simpleStore) GetAllTweets() []tweet {
	store.Lock()
	defer store.Unlock()
	tweets := make([]tweet, 0)

	for _, v := range store.tweetStorage {
		tweets = append(tweets, v)
	}

	return tweets
}

func (store *simpleStore) CreateTweet(tweet tweet) {
	store.Lock()
	defer store.Unlock()
	store.tweetStorage[tweet.ID] = tweet
}

func (store *simpleStore) DeleteTweet(id string) {
	store.Lock()
	defer store.Unlock()
	delete(store.tweetStorage, id)
}
