package main

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

type user struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func main() {
	r := mux.NewRouter()
	r.StrictSlash(true)
	r.HandleFunc("/api/user/", getAllUsers).Methods("GET")
	r.HandleFunc("/api/user/{id}/", getUser).Methods("GET")

	http.Handle("/", r)
	http.ListenAndServe(":8080", nil)
}

func getAllUsers(w http.ResponseWriter, r *http.Request) {
	respondData(w, r, user{"12345", "Test User"})
}

func getUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	respondData(w, r, user{vars["id"], "Test User"})
}

func respondData(w http.ResponseWriter, r *http.Request, data interface{}) error {
	return json.NewEncoder(w).Encode(data)
}
