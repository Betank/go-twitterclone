package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/SermoDigital/jose/jws"
	"github.com/gorilla/mux"
	nsq "github.com/nsqio/go-nsq"
)

type user struct {
	ID string `json:"id"`
}

type stats struct {
	Follow   int `json:"follow" bson:"follow"`
	Follower int `json:"follower" bson:"follower"`
	Tweets   int `json:"tweets" bson:"tweets"`
}

type tokenHandler struct {
	next http.Handler
}

type statsHandler struct{}

var store Storage
var createTweetConsumer, deleteTweetConsumer, newUserConsumer *nsq.Consumer
var nsqAddress string
var config *nsq.Config

func main() {
	store = NewMongoStorage()
	setupNSQ()

	router := mux.NewRouter()
	router.StrictSlash(true)
	router.Handle("/api/stats/", needsToken(&statsHandler{})).Methods("GET")
	router.HandleFunc("/api/stats/{userId}", statsForUser).Methods("GET")

	http.Handle("/", router)
	http.ListenAndServe(":8080", nil)
}

func (handler *statsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	respondData(w, r, store.GetStatsByUserID(r.Context().Value("userID").(string)))
}

func statsForUser(w http.ResponseWriter, r *http.Request) {
	respondData(w, r, store.GetStatsByUserID("12345"))
}

func (auth *tokenHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	jwt, err := jws.ParseJWTFromRequest(r)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	ctx := context.WithValue(r.Context(), "userID", jwt.Claims().Get("userID"))

	auth.next.ServeHTTP(w, r.WithContext(ctx))
}

func needsToken(handler http.Handler) http.Handler {
	return &tokenHandler{handler}
}

func respondData(w http.ResponseWriter, r *http.Request, data interface{}) error {
	return json.NewEncoder(w).Encode(data)
}

func setupNSQ() {
	nsqAddress = os.Getenv("NSQ_ADDRESS")
	if nsqAddress == "" {
		log.Fatal("nsq address not set")
	}

	config = nsq.NewConfig()

	createTweetConsumer = setupNSQConsumerHandler("create_tweet", updateTweetCount)
	deleteTweetConsumer = setupNSQConsumerHandler("delete_tweet", reduceTweetCount)
	newUserConsumer = setupNSQConsumerHandler("new_user", newUser)
}

func setupNSQConsumerHandler(topic string, handler func(message *nsq.Message) error) *nsq.Consumer {
	consumer, err := nsq.NewConsumer(topic, "stats", config)
	if err != nil {
		log.Fatal("error while creating producer")
	}
	consumer.AddHandler(nsq.HandlerFunc(handler))

	err = consumer.ConnectToNSQD(nsqAddress)
	if err != nil {
		log.Fatal("Could not connect to nsq")
	}

	return consumer
}

func updateTweetCount(message *nsq.Message) error {
	user, err := getUserFromMessage(message)
	if err != nil {
		log.Println("error while recieving message ", err.Error())
		return err
	}
	store.AddTweet(user.ID)
	return nil
}

func reduceTweetCount(message *nsq.Message) error {
	user, err := getUserFromMessage(message)
	if err != nil {
		log.Println("error while recieving message ", err.Error())
		return err
	}
	store.RemoveTweet(user.ID)
	return nil
}

func newUser(message *nsq.Message) error {
	user, err := getUserFromMessage(message)
	if err != nil {
		log.Println("error while recieving message ", err.Error())
		return err
	}
	store.CreateNewEntry(user.ID)
	return nil
}

func getUserFromMessage(message *nsq.Message) (user, error) {
	content := &struct {
		User user `json:"user"`
	}{}
	err := json.Unmarshal(message.Body, content)
	if err != nil {
		return user{}, err
	}

	return content.User, nil
}
