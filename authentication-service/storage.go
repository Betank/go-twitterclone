package main

import (
	"errors"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type User struct {
	Name     string `json:"username" bson:"username"`
	ID       string `json:"id" bson:"id"`
	Password string `json:"password" bson:"password"`
}

type Storage interface {
	AddUser(user User) error
	GetUserByName(name string) (User, error)
	RemoveUser(id string)
}

type mongoStorage struct {
	session *mgo.Session
}

var ErrUserNotFound = errors.New("User not found")
var ErrUserAlreadyExists = errors.New("User already exists")

func NewMongoStorage() *mongoStorage {
	session, err := mgo.Dial("mongo")
	if err != nil {

	}
	return &mongoStorage{session}
}

func (store *mongoStorage) AddUser(user User) error {
	sessionCopy := store.session.Copy()
	defer sessionCopy.Close()

	return sessionCopy.DB("gotwitterclone").C("user").Insert(&user)
}

func (store *mongoStorage) GetUserByName(name string) (User, error) {
	sessionCopy := store.session.Copy()
	defer sessionCopy.Close()

	var user User
	err := sessionCopy.DB("gotwitterclone").C("user").Find(bson.M{"username": name}).One(&user)
	if err != nil {
		return user, err
	}
	return user, nil
}

func (store *mongoStorage) RemoveUser(id string) {
	sessionCopy := store.session.Copy()
	defer sessionCopy.Close()

	sessionCopy.DB("gotwitterclone").C("user").Remove(bson.M{"id": id})
}
