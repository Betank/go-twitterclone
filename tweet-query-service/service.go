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
var createTweetConsumer *nsq.Consumer

func main() {
	store = &simpleStore{
		tweetStorage: make(map[string]tweet),
	}

	setupNSQConsumerHandler()

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
