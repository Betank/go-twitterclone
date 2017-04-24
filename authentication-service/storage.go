package main

import "sync"
import "errors"

type User struct {
	Name     string `json:"username"`
	ID       string `json:"id"`
	Password string `json:"password"`
}

type Storage interface {
	AddUser(user User) error
	GetUserByName(name string) (User, error)
	RemoveUser(id string)
}

type simpleStorage struct {
	sync.Mutex
	userStore map[string]User
}

var ErrUserNotFound = errors.New("User not found")
var ErrUserAlreadyExists = errors.New("User already exists")

func (store *simpleStorage) AddUser(user User) error {
	store.Lock()
	defer store.Unlock()
	if _, contains := store.userStore[user.Name]; !contains {
		store.userStore[user.Name] = user
		return nil
	}
	return ErrUserAlreadyExists
}

func (store *simpleStorage) GetUserByName(name string) (User, error) {
	store.Lock()
	defer store.Unlock()

	for _, v := range store.userStore {
		if v.Name == name {
			return v, nil
		}
	}
	return User{}, ErrUserNotFound
}

func (store *simpleStorage) RemoveUser(id string) {
	store.Lock()
	defer store.Unlock()

	delete(store.userStore, id)
}
