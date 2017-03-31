package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	nsq "github.com/nsqio/go-nsq"
)

var store Storage
var createTweetConsumer *nsq.Consumer

func main() {
	store = &simpleStore{
		tweetStorage: make(map[string]tweet),
	}

	setupNSQConsumerHandler()

	r := mux.NewRouter()
	r.StrictSlash(true)
	r.HandleFunc("/api/tweets/", allTweets).Methods("GET")
	r.HandleFunc("/api/tweet/{id}/", getTweet).Methods("GET")

	http.Handle("/", r)
	http.ListenAndServe(":8080", nil)
}

func allTweets(w http.ResponseWriter, r *http.Request) {
	respondData(w, r, store.GetAllTweets())
}

func getTweet(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	respondData(w, r, store.GetTweetById(vars["id"]))
}

func respondData(w http.ResponseWriter, r *http.Request, data interface{}) error {
	return json.NewEncoder(w).Encode(data)
}

func setupNSQConsumerHandler() {
	nsqAddress := os.Getenv("NSQ_ADDRESS")
	if nsqAddress == "" {
		log.Fatal("nsq address not set")
	}

	config := nsq.NewConfig()

	var err error
	createTweetConsumer, err = nsq.NewConsumer("create_tweet", "tweet", config)
	if err != nil {
		log.Fatal("error while creating producer")
	}
	createTweetConsumer.AddHandler(nsq.HandlerFunc(storeTweet))

	err = createTweetConsumer.ConnectToNSQD(nsqAddress)
	if err != nil {
		log.Fatal("Could not connect to nsq")
	}
}

func storeTweet(message *nsq.Message) error {
	tweet := &tweet{}
	err := json.Unmarshal(message.Body, tweet)
	if err != nil {
		log.Println("error while recieving message ", err.Error())
		return err
	}
	store.CreateTweet(*tweet)
	return nil
}
