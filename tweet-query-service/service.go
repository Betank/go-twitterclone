package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	nsq "github.com/nsqio/go-nsq"
)

var store Storage
var createTweetConsumer, deleteTweetConsumer *nsq.Consumer
var config *nsq.Config
var nsqAddress string

func main() {
	store = NewMongoStorage()

	setupNSQ()

	http.Handle("/", Router())
	http.ListenAndServe(":8080", nil)
}

func setupNSQ() {
	nsqAddress = os.Getenv("NSQ_ADDRESS")
	if nsqAddress == "" {
		log.Fatal("nsq address not set")
	}

	config = nsq.NewConfig()

	createTweetConsumer = setupNSQConsumerHandler("create_tweet", storeTweet)
	deleteTweetConsumer = setupNSQConsumerHandler("delete_tweet", deleteTweet)
}

func setupNSQConsumerHandler(topic string, handler func(message *nsq.Message) error) *nsq.Consumer {
	consumer, err := nsq.NewConsumer(topic, "tweet", config)
	if err != nil {
		log.Fatal("error while creating producer")
	}
	consumer.AddHandler(nsq.HandlerFunc(handler))

	err = consumer.ConnectToNSQD(nsqAddress)
	if err != nil {
		log.Fatal("Could not connect to nsq")
	}

	return consumer
}

func storeTweet(message *nsq.Message) error {
	tweet := &tweet{}
	err := json.Unmarshal(message.Body, tweet)
	if err != nil {
		log.Println("error while recieving message ", err.Error())
		return err
	}
	store.CreateTweet(*tweet)
	return nil
}

func deleteTweet(message *nsq.Message) error {
	tweetInfo := &struct {
		ID string `json:"id"`
	}{}
	err := json.Unmarshal(message.Body, tweetInfo)
	if err != nil {
		log.Println("error while recieving message ", err.Error())
		return err
	}
	store.DeleteTweet(tweetInfo.ID)
	return nil
}
