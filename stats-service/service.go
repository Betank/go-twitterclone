package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	nsq "github.com/nsqio/go-nsq"
)

type user struct {
	ID string `json:"id"`
}

type stats struct {
	Follow   int `json:"follow"`
	Follower int `json:"follower"`
	Tweets   int `json:"tweets"`
}

var store Storage
var createTweetConsumer *nsq.Consumer

func main() {
	store = simpleStoreMockUp()
	setupNSQConsumerHandler()

	router := mux.NewRouter()
	router.StrictSlash(true)
	router.HandleFunc("/api/stats/", statsForCurrentUser).Methods("GET")
	router.HandleFunc("/api/stats/{userId}", statsForUser).Methods("GET")

	http.Handle("/", router)
	http.ListenAndServe(":8080", nil)
}

func statsForCurrentUser(w http.ResponseWriter, r *http.Request) {
	respondData(w, r, store.GetStatsByUserID("12345"))
}

func statsForUser(w http.ResponseWriter, r *http.Request) {
	respondData(w, r, store.GetStatsByUserID("12345"))
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
	createTweetConsumer, err = nsq.NewConsumer("create_tweet", "stats", config)
	if err != nil {
		log.Fatal("error while creating producer")
	}
	createTweetConsumer.AddHandler(nsq.HandlerFunc(updateTweetCount))

	err = createTweetConsumer.ConnectToNSQD(nsqAddress)
	if err != nil {
		log.Fatal("Could not connect to nsq")
	}
}

func updateTweetCount(message *nsq.Message) error {
	content := &struct {
		User user `json:"user"`
	}{}
	err := json.Unmarshal(message.Body, content)
	if err != nil {
		log.Println("error while recieving message ", err.Error())
		return err
	}
	store.UpdateTweetCount(content.User.ID)
	return nil
}
