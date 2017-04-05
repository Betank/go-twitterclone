package main

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

type stats struct {
	Follow   int `json:"follow"`
	Follower int `json:"follower"`
	Tweets   int `json:"tweets"`
}

func main() {
	router := mux.NewRouter()
	router.StrictSlash(true)
	router.HandleFunc("/api/stats/", statsForCurrentUser).Methods("GET")
	router.HandleFunc("/api/stats/{userId}", statsForUser).Methods("GET")

	http.Handle("/", router)
	http.ListenAndServe(":8080", nil)
}

func statsForCurrentUser(w http.ResponseWriter, r *http.Request) {
	respondData(w, r, stats{1, 1, 1})
}

func statsForUser(w http.ResponseWriter, r *http.Request) {
	respondData(w, r, stats{0, 0, 0})
}

func respondData(w http.ResponseWriter, r *http.Request, data interface{}) error {
	return json.NewEncoder(w).Encode(data)
}
