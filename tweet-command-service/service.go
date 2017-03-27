package main

import (
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()
	r.StrictSlash(true)
	r.HandleFunc("/api/tweet/", createTweet).Methods("POST")
	r.HandleFunc("/api/tweet/{id}/", deleteTweet).Methods("DELETE")

	http.Handle("/", r)
	http.ListenAndServe(":8080", nil)
}

func createTweet(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func deleteTweet(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}
