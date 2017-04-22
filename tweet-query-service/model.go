package main

type tweet struct {
	ID   string `json:"id" bson:"id"`
	User user   `json:"user" bson:"user"`
	Text string `json:"text" bson:"text"`
}

type user struct {
	Name string `json:"name" bson:"name"`
	ID   string `json:"id" bson:"id"`
}
