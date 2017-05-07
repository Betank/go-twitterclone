package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/SermoDigital/jose/jws"
)

var testTweet = tweet{ID: "1234", Text: "tweet", User: user{"1", "test"}}

func TestGetTweetByIdHandler(t *testing.T) {
	store = &mockStorage{tweetById: testTweet}

	req, err := http.NewRequest("GET", "/api/tweet/1234/", nil)
	if err != nil {
		t.Error(err)
	}

	recorder := httptest.NewRecorder()

	Router().ServeHTTP(recorder, req)

	tweet := tweet{}
	err = json.NewDecoder(recorder.Body).Decode(&tweet)
	if err != nil {
		t.Error(err)
	}

	if tweet != testTweet {
		t.Errorf("expected %v but got %v", testTweet, tweet)
	}
}

func TestGetTweetByIdHandlerNotFound(t *testing.T) {
	store = &mockStorage{}

	req, err := http.NewRequest("GET", "/api/tweet/1234/", nil)
	if err != nil {
		t.Error(err)
	}

	recorder := httptest.NewRecorder()

	Router().ServeHTTP(recorder, req)

	if recorder.Code != http.StatusNotFound {
		t.Errorf("expected response code %d but got %d", http.StatusNotFound, recorder.Code)
	}
}

func TestGetTweetByUserId(t *testing.T) {
	store = &mockStorage{tweetsByUserId: []tweet{testTweet}, userId: testTweet.User.ID}

	jwt, err := createJWTForUser(testTweet.User)
	if err != nil {
		t.Error(err)
	}

	req, err := http.NewRequest("GET", "/api/tweets/user/", nil)
	if err != nil {
		t.Error(err)
	}
	req.Header.Set("Authorization", "bearer "+jwt)

	recorder := httptest.NewRecorder()

	Router().ServeHTTP(recorder, req)

	var tweets []tweet
	err = json.NewDecoder(recorder.Body).Decode(&tweets)
	if err != nil {
		t.Error(err)
	}

	if tweets[0] != testTweet {
		t.Errorf("expected %v but got %v", []tweet{testTweet}, tweets)
	}
}

func TestGetTweetByUserIdEmptyArray(t *testing.T) {
	store = &mockStorage{tweetsByUserId: []tweet{testTweet}, userId: "2"}

	jwt, err := createJWTForUser(testTweet.User)
	if err != nil {
		t.Error(err)
	}

	req, err := http.NewRequest("GET", "/api/tweets/user/", nil)
	if err != nil {
		t.Error(err)
	}
	req.Header.Set("Authorization", "bearer "+jwt)

	recorder := httptest.NewRecorder()

	Router().ServeHTTP(recorder, req)

	var tweets []tweet
	err = json.NewDecoder(recorder.Body).Decode(&tweets)
	if err != nil {
		t.Error(err)
	}

	if len(tweets) > 0 {
		t.Errorf("expected len %d but got len %d", 0, len(tweets))
	}
}

func TestGetTweetsByUserIdUnauthorized(t *testing.T) {
	store = &mockStorage{tweetsByUserId: []tweet{testTweet}, userId: "2"}
	req, err := http.NewRequest("GET", "/api/tweets/user/", nil)
	if err != nil {
		t.Error(err)
	}

	recorder := httptest.NewRecorder()

	Router().ServeHTTP(recorder, req)

	if recorder.Code != http.StatusUnauthorized {
		t.Errorf("expected response code %d but got %d", http.StatusNotFound, recorder.Code)
	}
}

func createJWTForUser(user user) (string, error) {
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

type mockStorage struct {
	tweetById      tweet
	tweetsByUserId []tweet
	userId         string
}

func (store *mockStorage) CreateTweet(tweet tweet) {
}

func (store *mockStorage) DeleteTweet(id string) {

}

func (store *mockStorage) GetTweetById(id string) (tweet, error) {
	if store.tweetById.ID != id {
		return tweet{}, errors.New("not found")
	}
	return store.tweetById, nil
}

func (store *mockStorage) GetTweetsByUserId(id string) []tweet {
	if store.userId == id {
		return store.tweetsByUserId
	}

	return make([]tweet, 0)
}
