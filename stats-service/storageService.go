package main

import (
	"sync"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type Storage interface {
	CreateNewEntry(id string)
	GetStatsByUserID(id string) stats
	AddTweet(id string)
	RemoveTweet(id string)
	UpdateFollowCount(id string)
	UpdateFollowerCount(id string)
	RemoveStats(id string)
}

type simpleStore struct {
	sync.Mutex
	statStore map[string]*stats
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

func (store *mongoStorage) CreateNewEntry(id string) {
	sessionCopy := store.session.Copy()
	defer sessionCopy.Close()

	entry := struct {
		id   string `bson:"id"`
		stat stats  `bson:"stats"`
	}{id, stats{}}

	sessionCopy.DB("gotwitterclone").C("stats").Insert(&entry)
}

func (store *mongoStorage) GetStatsByUserID(id string) stats {
	sessionCopy := store.session.Copy()
	defer sessionCopy.Close()

	entry := struct {
		id   string `bson:"id"`
		stat stats  `bson:"stats"`
	}{}
	sessionCopy.DB("gotwitterclone").C("stats").Find(bson.M{"id": id}).One(&entry)
	return entry.stat
}

func (store *mongoStorage) AddTweet(id string) {
	sessionCopy := store.session.Copy()
	defer sessionCopy.Close()

	entry := struct {
		id   string `bson:"id"`
		stat stats  `bson:"stats"`
	}{}
	sessionCopy.DB("gotwitterclone").C("stats").Find(bson.M{"id": id}).One(&entry)

	entry.stat.Tweets++
	sessionCopy.DB("gotwitterclone").C("stats").Update(bson.M{"id": id}, &entry)
}

func (store *mongoStorage) RemoveTweet(id string) {
	sessionCopy := store.session.Copy()
	defer sessionCopy.Close()

	entry := struct {
		id   string `bson:"id"`
		stat stats  `bson:"stats"`
	}{}
	sessionCopy.DB("gotwitterclone").C("stats").Find(bson.M{"id": id}).One(&entry)

	entry.stat.Tweets--
	sessionCopy.DB("gotwitterclone").C("stats").Update(bson.M{"id": id}, &entry)
}

func (store *mongoStorage) UpdateFollowCount(id string) {
	sessionCopy := store.session.Copy()
	defer sessionCopy.Close()

	entry := struct {
		id   string `bson:"id"`
		stat stats  `bson:"stats"`
	}{}
	sessionCopy.DB("gotwitterclone").C("stats").Find(bson.M{"id": id}).One(&entry)

	entry.stat.Follow++
	sessionCopy.DB("gotwitterclone").C("stats").Update(bson.M{"id": id}, &entry)
}

func (store *mongoStorage) UpdateFollowerCount(id string) {
	sessionCopy := store.session.Copy()
	defer sessionCopy.Close()

	entry := struct {
		id   string `bson:"id"`
		stat stats  `bson:"stats"`
	}{}
	sessionCopy.DB("gotwitterclone").C("stats").Find(bson.M{"id": id}).One(&entry)

	entry.stat.Follower++
	sessionCopy.DB("gotwitterclone").C("stats").Update(bson.M{"id": id}, &entry)
}

func (store *mongoStorage) RemoveStats(id string) {
	sessionCopy := store.session.Copy()
	defer sessionCopy.Close()
	sessionCopy.DB("gotwitterclone").C("stats").Remove(bson.M{"id": id})

}
