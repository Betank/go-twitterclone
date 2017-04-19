package main

import (
	"encoding/json"
	"net/http"

	"io/ioutil"

	"errors"

	"github.com/SermoDigital/jose/jws"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type tokenResponse struct {
	Token string `json:"token"`
}

var storage Storage

func main() {
	storage = &simpleStorage{
		userStore: make(map[string]User),
	}

	r := mux.NewRouter()
	r.HandleFunc("/api/login/", login)
	r.HandleFunc("/api/register/", register)

	http.Handle("/", r)
	http.ListenAndServe(":8080", nil)
}

func login(w http.ResponseWriter, r *http.Request) {
	username, password, err := parseRequest(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	user, err := checkCredentials(username, password)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	jwt, err := generateJWT(user)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	token := &tokenResponse{jwt}

	respondData(w, r, token)
}

func parseRequest(r *http.Request) (string, string, error) {
	err := r.ParseForm()
	if err != nil {
		return "", "", err
	}
	username := r.FormValue("username")
	password := r.FormValue("password")

	return username, password, nil
}

func checkCredentials(username, password string) (User, error) {
	user, err := storage.GetUserByName(username)
	if err != nil {
		return User{}, err
	}
	if user.Password != password {
		return User{}, errors.New("wrong password")
	}
	return user, nil
}

func register(w http.ResponseWriter, r *http.Request) {
	user := &User{}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = json.Unmarshal(body, user)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	user.ID = uuid.New().String()
	err = storage.AddUser(*user)
	if err != nil {
		w.WriteHeader(http.StatusConflict)
		return
	}
}

func generateJWT(user User) (string, error) {
	payload := make(jws.Claims)
	payload.Set("userID", user.ID)
	payload.Set("userName", user.Name)

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
