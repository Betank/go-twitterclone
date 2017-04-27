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

	id, err := createTweet(user, "hello tweet")
	if err != nil {
		t.Error(err)
	}
	_, err = awaitTweet(user, id)
	if err != nil {
		t.Error(err)
	}
	id2, err := createTweet(user, "hello tweet2")
	if err != nil {
		t.Error(err)
	}
	_, err = awaitTweet(user, id2)
	if err != nil {
		t.Error(err)
	}

	tweets, err := getTweetsForUser(user)
	if err != nil {
		t.Error(err)
	}

	if len(tweets) != 2 {
		t.Error("Should have %d tweets, but has %d tweets", 2, len(tweets))
	}
}

func assertTweetEqual(gotTweet, wantTweet tweet, t *testing.T) {
	if gotTweet != wantTweet {
		t.Errorf("tweets are not equal: got %v want %v", gotTweet, wantTweet)
	}
}
