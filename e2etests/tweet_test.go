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

func TestGetTweetsForCurrentUser(t *testing.T) {
	user := user{ID: "0078142e-a2cd-4755-8167-da5cf856294a", Name: "user2"}

	tweets, err := createMultipleTweetsAndAwait(user, "test1", "test2")
	if err != nil {
		t.Error(err)
	}

	userTweets, err := getTweetsForUser(user)
	if err != nil {
		t.Error(err)
	}

	assertTweetsEqual(userTweets, tweets, t)
}

func assertTweetEqual(gotTweet, wantTweet tweet, t *testing.T) {
	if gotTweet != wantTweet {
		t.Errorf("tweets are not equal: got %v want %v", gotTweet, wantTweet)
	}
}

func assertTweetsEqual(gotTweets, wantTweets []tweet, t *testing.T) {
	for _, tweet := range gotTweets {
		if !containsTweet(wantTweets, tweet) {
			t.Errorf("tweet %v isn't part of %v", tweet, wantTweets)
		}
	}
}

func containsTweet(tweets []tweet, tweet tweet) bool {
	for _, currTweet := range tweets {
		if currTweet == tweet {
			return true
		}
	}
	return false
}
