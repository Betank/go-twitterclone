package e2etests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/url"

	mgo "gopkg.in/mgo.v2"

	"time"

	"io"

	"os"

	"fmt"

	"errors"

	"github.com/SermoDigital/jose/jws"
)

var gatewayURL = "http://localhost"
var mongoDbAddr = "localhost:27017"

var loginFailedError = errors.New("Login failed")
var userConflictError = errors.New("User already exists")

func setGatewayURL() {
	url := os.Getenv("GATEWAY_URL")
	if url != "" {
		gatewayURL = url
	}
}

func setMongoDbAddr() {
	address := os.Getenv("MONGO_ADDRESS")
	if address != "" {
		mongoDbAddr = address
	}
}

func dropDB() {
	session, err := mgo.Dial(mongoDbAddr)
	if err != nil {
		fmt.Printf("unable to drop database because of %s\n", err)
		return
	}
	err = session.DB("gotwitterclone").DropDatabase()
	if err != nil {
		fmt.Printf("unable to drop database because of %s\n", err)
	}
}

func createMultipleTweetsAndAwait(user user, text ...string) ([]tweet, error) {
	tweets := make([]tweet, 0)
	for _, content := range text {
		id, err := createTweet(user, content)
		if err != nil {
			return tweets, err
		}
		tweet, err := awaitTweet(user, id)
		if err != nil {
			return tweets, err
		}
		tweets = append(tweets, tweet)
	}
	return tweets, nil
}

func createTweet(user user, text string) (string, error) {
	reqBody := []byte(text)
	req, err := createNewAuthHeaderRequest(user, "POST", gatewayURL+"/api/tweet/", bytes.NewReader(reqBody))
	if err != nil {
		return "", err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}

	respObject := &struct {
		ID string `json:"id"`
	}{}
	err = json.NewDecoder(resp.Body).Decode(respObject)
	if err != nil {
		return "", err
	}
	return respObject.ID, nil
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

func awaitTweet(user user, tweetID string) (tweet, error) {
	tweetRequest, err := createNewAuthHeaderRequest(user,
		"GET",
		gatewayURL+"/api/tweet/"+tweetID+"/",
		nil)
	if err != nil {
		return tweet{}, err
	}

	resp := doRequestUntilSuccess(tweetRequest, http.StatusOK, 200)

	tweetResponseBody := &tweet{}
	err = json.NewDecoder(resp.Body).Decode(tweetResponseBody)
	if err != nil {
		return tweet{}, err
	}
	return *tweetResponseBody, nil
}

func getTweetsForUser(user user) ([]tweet, error) {
	tweets := make([]tweet, 0)
	tweetsRequest, err := createNewAuthHeaderRequest(user,
		"GET",
		gatewayURL+"/api/tweets/user/",
		nil)
	if err != nil {
		return tweets, err
	}
	resp, err := http.DefaultClient.Do(tweetsRequest)
	if err != nil {
		return tweets, err
	}
	err = json.NewDecoder(resp.Body).Decode(&tweets)
	if err != nil {
		return tweets, err
	}
	return tweets, nil
}

func deleteTweet(user user, tweetID string) error {
	deleteTweetRequest, err := createNewAuthHeaderRequest(user,
		"DELETE",
		gatewayURL+"/api/tweet/"+tweetID+"/",
		nil)
	if err != nil {
		return err
	}
	resp, err := http.DefaultClient.Do(deleteTweetRequest)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("request failed with status: %d", resp.StatusCode)
	}
	return nil
}

func awaitTweetDeleted(user user, tweetID string) error {
	tweetRequest, err := createNewAuthHeaderRequest(user,
		"GET",
		gatewayURL+"/api/tweet/"+tweetID+"/",
		nil)
	if err != nil {
		return err
	}

	resp := doRequestUntilSuccess(tweetRequest, http.StatusNotFound, 200)

	if resp.StatusCode != 404 {
		return fmt.Errorf("wrong status %d", resp.StatusCode)
	}

	return nil
}

func awaitStats(username, password string) (stats, error) {
	token, err := userLogin(username, password)
	if err != nil {
		return stats{}, err
	}

	statsRequest, err := http.NewRequest(
		"GET",
		gatewayURL+"/api/stats/",
		nil)
	if err != nil {
		return stats{}, err
	}

	statsRequest.Header.Set("Authorization", "bearer "+token)

	resp := doRequestUntilSuccess(statsRequest, http.StatusOK, 200)
	if resp.StatusCode != 200 {
		return stats{}, fmt.Errorf("wrong status %d", resp.StatusCode)
	}

	statsBody := stats{}
	err = json.NewDecoder(resp.Body).Decode(&statsBody)
	if err != nil {
		return stats{}, err
	}
	return statsBody, nil
}

func createNewAuthHeaderRequest(user user, method, url string, body io.Reader) (*http.Request, error) {
	var request *http.Request
	jwt, err := createJWTForUser(user)
	if err != nil {
		return request, err
	}

	request, err = http.NewRequest(method, url, body)
	if err != nil {
		return request, err
	}
	request.Header.Add("Authorization", "bearer "+jwt)
	return request, nil
}

func doRequestUntilSuccess(r *http.Request, wantStatus int, repeats int) *http.Response {
	var count int
	ticker := time.NewTicker(100 * time.Millisecond)
	stop := make(chan *http.Response, 1)

	go func() {
		for {
			select {
			case <-ticker.C:
				resp, err := http.DefaultClient.Do(r)
				if (err == nil && resp.StatusCode == wantStatus) || repeats == count {
					stop <- resp
				}
				count++
			case <-stop:
				return
			}
		}
	}()

	return <-stop
}

func registerNewUser(name, password string) error {
	user := struct {
		Name     string `json:"username"`
		Password string `json:"password"`
	}{name, password}

	body, err := json.Marshal(&user)
	if err != nil {
		return err
	}

	resp, err := http.Post(gatewayURL+"/api/register/", "application/json", bytes.NewReader(body))
	if err != nil {
		return err
	}

	if resp.StatusCode == 409 {
		return userConflictError
	}

	if resp.StatusCode != 200 {
		return errors.New("error while creating user")
	}
	return nil
}

func userLogin(name, password string) (string, error) {
	resp, err := http.PostForm(gatewayURL+"/api/login/", url.Values{"username": {name}, "password": {password}})
	if err != nil {
		return "", err
	}

	if resp.StatusCode == 401 {
		return "", loginFailedError
	}

	if resp.StatusCode != 200 {
		return "", errors.New("error while creating user")
	}

	jwt := &struct {
		Token string `json:"token"`
	}{}
	err = json.NewDecoder(resp.Body).Decode(jwt)
	if err != nil {
		return "", nil
	}
	return jwt.Token, nil
}
