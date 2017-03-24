package main

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/api/tweets", allTweets).Methods("GET")
	r.HandleFunc("/api/tweet/{id}", getTweet).Methods("GET")

	http.Handle("/", r)
	http.ListenAndServe(":8080", nil)
}

func allTweets(w http.ResponseWriter, r *http.Request) {
	mockedTweet1 := tweet{
		ID:   "1ABC",
		User: user{"test"},
		Text: "test tweet 1",
	}

	mockedTweet2 := tweet{
		ID:   "1ABD",
		User: user{"test"},
		Text: "test tweet 2",
	}

	tweets, err := json.Marshal([]tweet{mockedTweet1, mockedTweet2})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Write(tweets)
}

func getTweet(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	mockedTweet1 := tweet{
		ID:   vars["id"],
		User: user{"test"},
		Text: "test tweet 1",
	}

	tweet, err := json.Marshal(mockedTweet1)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Write(tweet)
}
