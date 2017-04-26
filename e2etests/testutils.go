package e2etests

import (
	"bytes"
	"encoding/json"
	"net/http"

	"time"

	"io"

	"os"

	"github.com/SermoDigital/jose/jws"
)

var gatewayURL = "http://localhost"

func setGatewayURL() {
	url := os.Getenv("GATEWAY_URL")
	if url != "" {
		gatewayURL = url
	}
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
