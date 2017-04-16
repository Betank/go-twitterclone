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

type authHandler struct {
	next http.Handler
}

var tweetCommandServiceURL *url.URL
var tweetQueryServiceURL *url.URL
var userCommandServiceURL *url.URL
var userQueryServiceURL *url.URL
var statServiceURL *url.URL
var authServiceURL *url.URL

var dir = flag.String("d", "./client/public", "client location")

func main() {

	flag.Parse()

	setUpServiceURLs()

	r := mux.NewRouter()
	r.StrictSlash(true)
	r.Handle("/api/tweet/", mustAuth(proxy(tweetCommandServiceURL))).Methods("POST", "DELETE")
	r.Handle("/api/tweet/{id}/", mustAuth(proxy(tweetCommandServiceURL))).Methods("POST", "DELETE")

	r.Handle("/api/tweets/", mustAuth(proxy(tweetQueryServiceURL))).Methods("GET")
	r.Handle("/api/tweet/{id}/", mustAuth(proxy(tweetQueryServiceURL))).Methods("GET")
	r.Handle("/api/tweets/user/", mustAuth(proxy(tweetQueryServiceURL))).Methods("GET")

	r.Handle("/api/user/", mustAuth(proxy(userCommandServiceURL))).Methods("POST", "DELETE")
	r.Handle("/api/user/{id}/", proxy(userCommandServiceURL)).Methods("POST", "DELETE")

	r.Handle("/api/user/", mustAuth(proxy(userQueryServiceURL))).Methods("GET")
	r.Handle("/api/user/{id}/", mustAuth(proxy(userQueryServiceURL))).Methods("GET")

	r.Handle("/api/stats/", mustAuth(proxy(statServiceURL))).Methods("GET")
	r.Handle("/api/stats/{userId}", mustAuth(proxy(statServiceURL))).Methods("GET")

	r.Handle("/api/login/", proxy(authServiceURL)).Methods("POST")

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
	statServiceURL = createURL("STATS_SERVICE_URL")
	authServiceURL = createURL("AUTH_SERVICE_URL")
}

func createURL(env string) *url.URL {
	url, err := url.Parse(os.Getenv(env))
	if err != nil {
		log.Fatal("Error while creating endpoints: ", err.Error())
	}

	return url
}

func mustAuth(handler http.Handler) http.Handler {
	return &authHandler{handler}
}

func (handler *authHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("Authorization")
	if token == "" {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	handler.next.ServeHTTP(w, r)
}
