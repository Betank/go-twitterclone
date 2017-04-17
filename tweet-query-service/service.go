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

type authHandler struct {
	next http.Handler
}
type tweetHandler struct{}

var store Storage
var createTweetConsumer, deleteTweetConsumer *nsq.Consumer
var config *nsq.Config
var nsqAddress string

func main() {
	store = &simpleStore{
		tweetStorage: make(map[string]tweet),
	}

	setupNSQ()

	r := mux.NewRouter()
	r.StrictSlash(true)
	r.HandleFunc("/api/tweet/{id}/", getTweet).Methods("GET")
	r.Handle("/api/tweets/user/", mustAuth(&tweetHandler{})).Methods("GET")
	r.HandleFunc("/api/tweets/", allTweets).Methods("GET")

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

func (handler *tweetHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	tweets := store.GetTweetsByUserId(r.Context().Value("userID").(string))
	respondData(w, r, tweets)
}

func (auth *authHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	jwt, err := jws.ParseJWTFromRequest(r)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	ctx := context.WithValue(r.Context(), "userID", jwt.Claims().Get("userID"))

	auth.next.ServeHTTP(w, r.WithContext(ctx))
}

func respondData(w http.ResponseWriter, r *http.Request, data interface{}) error {
	return json.NewEncoder(w).Encode(data)
}

func mustAuth(handler http.Handler) http.Handler {
	return &authHandler{handler}
}

func setupNSQ() {
	nsqAddress = os.Getenv("NSQ_ADDRESS")
	if nsqAddress == "" {
		log.Fatal("nsq address not set")
	}

	config = nsq.NewConfig()

	createTweetConsumer = setupNSQConsumerHandler("create_tweet", storeTweet)
	deleteTweetConsumer = setupNSQConsumerHandler("delete_tweet", deleteTweet)
}

func setupNSQConsumerHandler(topic string, handler func(message *nsq.Message) error) *nsq.Consumer {
	consumer, err := nsq.NewConsumer(topic, "tweet", config)
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

func deleteTweet(message *nsq.Message) error {
	tweetInfo := &struct {
		ID string `json:"id"`
	}{}
	err := json.Unmarshal(message.Body, tweetInfo)
	if err != nil {
		log.Println("error while recieving message ", err.Error())
		return err
	}
	store.DeleteTweet(tweetInfo.ID)
	return nil
}
