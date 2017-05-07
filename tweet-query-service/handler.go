package main

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/SermoDigital/jose/jws"
	"github.com/gorilla/mux"
)

type authHandler struct {
	next http.Handler
}
type tweetHandler struct{}

func Router() *mux.Router {
	r := mux.NewRouter()
	r.StrictSlash(true)
	r.HandleFunc("/api/tweet/{id}/", getTweet).Methods("GET")
	r.Handle("/api/tweets/user/", mustAuth(&tweetHandler{})).Methods("GET")

	return r
}

func getTweet(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	tweet, err := store.GetTweetById(vars["id"])
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	respondData(w, r, tweet)
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
