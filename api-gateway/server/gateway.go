package main

import (
	"flag"
	"log"
	"net/http"
	"net/url"
	"os"

	"net/http/httputil"

	"github.com/gorilla/mux"
)

var tweetCommandServiceURL *url.URL
var tweetQueryServiceURL *url.URL
var userCommandServiceURL *url.URL
var userQueryServiceURL *url.URL

var dir = flag.String("d", "./client/public", "client location")

func main() {

	flag.Parse()

	setUpServiceURLs()

	r := mux.NewRouter()
	r.StrictSlash(true)
	r.Handle("/api/tweet/", proxy(tweetCommandServiceURL)).Methods("POST", "DELETE")
	r.Handle("/api/tweet/{id}/", proxy(tweetCommandServiceURL)).Methods("POST", "DELETE")

	r.Handle("/api/tweets/", proxy(tweetQueryServiceURL)).Methods("GET")
	r.Handle("/api/tweet/{id}/", proxy(tweetQueryServiceURL)).Methods("GET")

	r.Handle("/api/user/", proxy(userCommandServiceURL)).Methods("POST", "DELETE")
	r.Handle("/api/user/{id}/", proxy(userCommandServiceURL)).Methods("POST", "DELETE")

	r.Handle("/api/user/", proxy(userQueryServiceURL)).Methods("GET")
	r.Handle("/api/user/{id}/", proxy(userQueryServiceURL)).Methods("GET")
	r.Handle("/api/stats/", proxy(userQueryServiceURL)).Methods("GET")

	r.Handle("/", http.FileServer(http.Dir(*dir)))
	r.PathPrefix("/dist/").Handler(http.FileServer(http.Dir(*dir)))

	http.Handle("/", r)
	http.ListenAndServe(":8080", nil)
}

func proxy(URL *url.URL) http.Handler {
	return httputil.NewSingleHostReverseProxy(URL)
}

func setUpServiceURLs() {
	tweetCommandServiceURL = createURL("TWEET_COMMAND_SERVICE_URL")
	tweetQueryServiceURL = createURL("TWEET_QUERY_SERVICE_URL")
	userCommandServiceURL = createURL("USER_COMMAND_SERVICE_URL")
	userQueryServiceURL = createURL("USER_QUERY_SERVICE_URL")
}

func createURL(env string) *url.URL {
	url, err := url.Parse(os.Getenv(env))
	if err != nil {
		log.Fatal("Error while creating endpoints: ", err.Error())
	}

	return url
}
