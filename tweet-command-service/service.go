package main

import (
	"net/http"
	"os"

	"log"

	"io/ioutil"

	"encoding/json"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	nsq "github.com/nsqio/go-nsq"
)

var producer *nsq.Producer

func main() {
	setupNSQProducer()

	r := mux.NewRouter()
	r.StrictSlash(true)
	r.HandleFunc("/api/tweet/", createTweet).Methods("POST")
	r.HandleFunc("/api/tweet/{id}/", deleteTweet).Methods("DELETE")

	http.Handle("/", r)
	http.ListenAndServe(":8080", nil)
}

func createTweet(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	id := uuid.New().String()
	tweet := tweet{
		ID:   id,
		Text: string(body),
		User: user{"Test user"},
	}

	event, err := json.Marshal(tweet)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = producer.Publish("create_tweet", event)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func deleteTweet(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func setupNSQProducer() {
	nsqAddress := os.Getenv("NSQ_ADDRESS")
	if nsqAddress == "" {
		log.Fatal("nsq address not set")
	}

	config := nsq.NewConfig()

	var err error
	producer, err = nsq.NewProducer(nsqAddress, config)
	if err != nil {
		log.Fatal("error while creating producer")
	}
}
