package main

import (
	"encoding/json"
	"net/http"

	"github.com/SermoDigital/jose/jws"
	"github.com/gorilla/mux"
)

type tokenResponse struct {
	Token string `json:"token"`
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/api/login/", login)

	http.Handle("/", r)
	http.ListenAndServe(":8080", nil)
}

func login(w http.ResponseWriter, r *http.Request) {
	jwt, err := generateJWT()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	token := &tokenResponse{jwt}

	respondData(w, r, token)
}

func generateJWT() (string, error) {
	payload := make(jws.Claims)
	payload.Set("userID", "12345")
	payload.Set("userName", "Test User")

	token := jws.NewJWT(payload, jws.GetSigningMethod("HS256"))

	jwt, err := token.Serialize([]byte("secret"))
	if err != nil {
		return "", err
	}

	return string(jwt), nil
}

func respondData(w http.ResponseWriter, r *http.Request, data interface{}) error {
	return json.NewEncoder(w).Encode(data)
}
