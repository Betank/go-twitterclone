package main

import (
	"net/http"
	"testing"

	"net/http/httptest"

	"encoding/json"

	"github.com/SermoDigital/jose/jws"
)

var testStats = stats{1, 1, 1}

func TestGetUserStats(t *testing.T) {
	store = &mockStorage{map[string]stats{"1": testStats}}
	jwt, err := createJWTForUser(user{"1"})
	if err != nil {
		t.Error(err)
	}

	req, err := http.NewRequest("GET", "/api/stats/", nil)
	if err != nil {
		t.Error(err)
	}
	req.Header.Set("Authorization", "bearer "+jwt)

	recorder := httptest.NewRecorder()
	Router().ServeHTTP(recorder, req)

	respStats := stats{}
	err = json.NewDecoder(recorder.Body).Decode(&respStats)
	if err != nil {
		t.Error(err)
	}

	if respStats != testStats {
		t.Errorf("stats are not equal, got %v and want %v", respStats, testStats)
	}
}

func TestGetStatsByUserId(t *testing.T) {
	var userStats = stats{2, 2, 2}
	store = &mockStorage{map[string]stats{"2": userStats}}

	req, err := http.NewRequest("GET", "/api/stats/2/", nil)
	if err != nil {
		t.Error(err)
	}

	recorder := httptest.NewRecorder()
	Router().ServeHTTP(recorder, req)

	respStats := stats{}
	err = json.NewDecoder(recorder.Body).Decode(&respStats)
	if err != nil {
		t.Error(err)
	}

	if respStats != userStats {
		t.Errorf("stats are not equal, got %v and want %v", respStats, testStats)
	}
}

func createJWTForUser(user user) (string, error) {
	payload := make(jws.Claims)
	payload.Set("userID", user.ID)

	token := jws.NewJWT(payload, jws.GetSigningMethod("HS256"))

	jwt, err := token.Serialize([]byte("secret"))
	if err != nil {
		return "", err
	}

	return string(jwt), nil
}

type mockStorage struct {
	userStats map[string]stats
}

func (store *mockStorage) CreateNewEntry(id string) {
}

func (store *mockStorage) GetStatsByUserID(id string) stats {
	return store.userStats[id]
}

func (store *mockStorage) AddTweet(id string) {
}

func (store *mockStorage) RemoveTweet(id string) {
}

func (store *mockStorage) UpdateFollowCount(id string) {
}

func (store *mockStorage) UpdateFollowerCount(id string) {
}

func (store *mockStorage) RemoveStats(id string) {

}
