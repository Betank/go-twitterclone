package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	nsq "github.com/nsqio/go-nsq"
)

var producer *nsq.Producer

func main() {
	setupNSQProducer()

	r := mux.NewRouter()
	r.StrictSlash(true)
	r.HandleFunc("/api/user/", createUser).Methods("POST")
	r.HandleFunc("/api/user/{id}", deleteUser).Methods("DELETE")

	http.Handle("/", r)
	http.ListenAndServe(":8080", nil)
}

func createUser(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func deleteUser(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
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
