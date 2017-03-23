package main

import (
	"encoding/json"
	"net/http"
)

func main() {
	http.HandleFunc("/api/tweets", allTweets)
	http.ListenAndServe(":8080", nil)
}

func allTweets(w http.ResponseWriter, r *http.Request) {
	mockedTweet1 := tweet{
		ID:   "1ABC",
		User: user{"test"},
		Text: "test tweet 1",
	}

	mockedTweet2 := tweet{
		ID:   "1ABC",
		User: user{"test"},
		Text: "test tweet 1",
	}

	tweets, err := json.Marshal([]tweet{mockedTweet1, mockedTweet2})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Write(tweets)
}
