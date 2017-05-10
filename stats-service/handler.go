package main

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/SermoDigital/jose/jws"
	"github.com/gorilla/mux"
)

type tokenHandler struct {
	next http.Handler
}

type statsHandler struct{}

func Router() *mux.Router {
	router := mux.NewRouter()
	router.StrictSlash(true)
	router.Handle("/api/stats/", needsToken(&statsHandler{})).Methods("GET")
	router.HandleFunc("/api/stats/{userId}/", statsForUser).Methods("GET")

	return router
}

func (handler *statsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	respondData(w, r, store.GetStatsByUserID(r.Context().Value("userID").(string)))
}

func statsForUser(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["userId"]
	respondData(w, r, store.GetStatsByUserID(id))
}

func (auth *tokenHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	jwt, err := jws.ParseJWTFromRequest(r)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	ctx := context.WithValue(r.Context(), "userID", jwt.Claims().Get("userID"))

	auth.next.ServeHTTP(w, r.WithContext(ctx))
}

func needsToken(handler http.Handler) http.Handler {
	return &tokenHandler{handler}
}

func respondData(w http.ResponseWriter, r *http.Request, data interface{}) error {
	return json.NewEncoder(w).Encode(data)
}
