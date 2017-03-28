package main

type tweet struct {
	ID   string `json:"id"`
	User user   `json:"user"`
	Text string `json:"text"`
}

type user struct {
	Name string `json:"name"`
}
