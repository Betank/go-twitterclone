package main

import (
	"errors"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type Storage interface {
	GetTweetById(id string) (tweet, error)
	GetTweetsByUserId(id string) []tweet
	CreateTweet(tweet tweet)
	DeleteTweet(id string)
}

type mongoStorage struct {
	session *mgo.Session
}

func NewMongoStorage() *mongoStorage {
	session, err := mgo.Dial("mongo")
	if err != nil {

	}
	return &mongoStorage{session}
}

func (store *mongoStorage) CreateTweet(tweet tweet) {
	sessionCopy := store.session.Copy()
	defer sessionCopy.Close()
	sessionCopy.DB("gotwitterclone").C("tweets").Insert(&tweet)
}

func (store *mongoStorage) DeleteTweet(id string) {
	sessionCopy := store.session.Copy()
	defer sessionCopy.Close()
	sessionCopy.DB("gotwitterclone").C("tweets").Remove(bson.M{"id": id})
}

func (store *mongoStorage) GetTweetById(id string) (tweet, error) {
	sessionCopy := store.session.Copy()
	defer sessionCopy.Close()

	respondTweet := tweet{}
	sessionCopy.DB("gotwitterclone").C("tweets").Find(bson.M{"id": id}).One(&respondTweet)
	if respondTweet == (tweet{}) {
		return respondTweet, errors.New("tweet not found")
	}

	return respondTweet, nil
}

func (store *mongoStorage) GetTweetsByUserId(id string) []tweet {
	sessionCopy := store.session.Copy()
	defer sessionCopy.Close()

	var tweets []tweet
	sessionCopy.DB("gotwitterclone").C("tweets").Find(bson.M{"user.id": id}).All(&tweets)
	return tweets
}
