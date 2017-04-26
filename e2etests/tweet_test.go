package e2etests

import (
	"os"
	"testing"
)

var testUser = user{ID: "6c9ce302-8de9-44fd-8161-05dc06925ad6", Name: "user"}

func TestMain(m *testing.M) {
	setGatewayURL()
	os.Exit(m.Run())
}

func TestShouldCreateTweet(t *testing.T) {
	tweetID, err := createTweet(testUser, "hello tweet")
	if err != nil {
		t.Error(err)
	}

	respondedTweet, err := awaitTweet(testUser, tweetID)
	if err != nil {
		t.Error(err)
	}

	assertTweetEqual(respondedTweet, tweet{
		ID:   tweetID,
		Text: "hello tweet",
		User: testUser,
	}, t)
}

func assertTweetEqual(gotTweet, wantTweet tweet, t *testing.T) {
	if gotTweet != wantTweet {
		t.Errorf("tweets are not equal: got %v want %v", gotTweet, wantTweet)
	}
}
