package main

import "sync"

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

func (store *simpleStore) CreateNewEntry(id string) {
	store.Lock()
	defer store.Unlock()
	store.statStore[id] = &stats{}
}

func (store *simpleStore) GetStatsByUserID(id string) stats {
	store.Lock()
	defer store.Unlock()
	stats := store.statStore[id]
	return *stats
}

func (store *simpleStore) AddTweet(id string) {
	store.Lock()
	defer store.Unlock()
	store.statStore[id].Tweets++
}

func (store *simpleStore) RemoveTweet(id string) {
	store.Lock()
	defer store.Unlock()
	if store.statStore[id].Tweets > 0 {
		store.statStore[id].Tweets--
	}
}

func (store *simpleStore) UpdateFollowCount(id string) {
	store.Lock()
	defer store.Unlock()
	store.statStore[id].Follow++
}

func (store *simpleStore) UpdateFollowerCount(id string) {
	store.Lock()
	defer store.Unlock()
	store.statStore[id].Follower++
}

func (store *simpleStore) RemoveStats(id string) {
	store.Lock()
	defer store.Unlock()
	delete(store.statStore, id)
}

func simpleStoreMockUp() *simpleStore {
	store := &simpleStore{
		statStore: make(map[string]*stats),
	}

	store.statStore["12345"] = &stats{}

	return store
}
