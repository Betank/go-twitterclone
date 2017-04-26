package main

import (
	"net/http"
	"os"

	"log"

	"io/ioutil"

	"encoding/json"

	"context"

	"github.com/SermoDigital/jose/jws"
	"github.com/SermoDigital/jose/jwt"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	nsq "github.com/nsqio/go-nsq"
)

type authHandler struct {
	next http.Handler
}

type tweetCreateHandler struct{}
type tweetDeleteHandler struct{}

var producer *nsq.Producer

func main() {
	setupNSQProducer()

	r := mux.NewRouter()
	r.StrictSlash(true)
	r.Handle("/api/tweet/", mustAuth(&tweetCreateHandler{})).Methods("POST")
	r.Handle("/api/tweet/{id}/", mustAuth(&tweetDeleteHandler{})).Methods("DELETE")

	http.Handle("/", r)
	http.ListenAndServe(":8080", nil)
}

func (handler *tweetCreateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	id := uuid.New().String()
	tweet := tweet{
		ID:   id,
		Text: string(body),
		User: r.Context().Value("user").(user),
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

	response := struct {
		ID string `json:"id"`
	}{id}
	if err = json.NewEncoder(w).Encode(response); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (handler *tweetDeleteHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	tweet := tweet{
		ID:   id,
		User: r.Context().Value("user").(user),
	}

	event, err := json.Marshal(tweet)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = producer.Publish("delete_tweet", event)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func mustAuth(handler http.Handler) http.Handler {
	return &authHandler{handler}
}

func (auth *authHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	jwt, err := jws.ParseJWTFromRequest(r)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	ctx := context.WithValue(r.Context(), "user", getUserFromJWT(jwt))

	auth.next.ServeHTTP(w, r.WithContext(ctx))
}

func getUserFromJWT(jwt jwt.JWT) user {
	return user{
		jwt.Claims().Get("userName").(string),
		jwt.Claims().Get("userID").(string),
	}
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
